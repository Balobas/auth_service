package config

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type ServiceConfig struct {
	model     *serviceConfigModel
	configEnv configEnv
}

func NewServiceConfig() *ServiceConfig {
	envCfg := configEnv{}
	ParseEnv(&envCfg)
	return &ServiceConfig{
		configEnv: envCfg,
	}
}

type serviceConfigModel struct {
	mu *sync.RWMutex

	minPasswordLen              int           `setting_name:"min_password_len" default:"6"`
	accessJwtTTL                time.Duration `setting_name:"access_jwt_ttl" default:"1h"`
	refreshJwtTTL               time.Duration `setting_name:"refresh_jwt_ttl" default:"24h"`
	verificationTokenLen        int64         `setting_name:"verification_token_len" default:"16"`
	sendVerificationInterval    time.Duration `setting_name:"send_verification_interval" default:"3m"`
	verificationWorkerBatchSize uint64        `setting_name:"verification_worker_batch_size" default:"10"`
	emailVerificationTemplate   string        `setting_name:"email_verification_template" default:"{{Подтвердите вашу почту перейдя по ссылке .Scheme/.Token }}"`
	httpVerificationScheme      string        `setting_name:"http_verification_scheme"`
}

func (c *ServiceConfig) MinPasswordLen() int {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.minPasswordLen
}

func (c *ServiceConfig) AccessJwtTTL() time.Duration {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.accessJwtTTL
}

func (c *ServiceConfig) RefreshJwtTTL() time.Duration {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.refreshJwtTTL
}

func (c *ServiceConfig) VerificationTokenLen() int64 {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.verificationTokenLen
}

func (c *ServiceConfig) SendVerificationInterval() time.Duration {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.sendVerificationInterval
}

func (c *ServiceConfig) VerificationWorkerBatchSize() uint64 {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.verificationWorkerBatchSize
}

func (c *ServiceConfig) EmailVerificationTemplate() string {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.emailVerificationTemplate
}

func (c *ServiceConfig) HttpVerificationScheme() string {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()
	return c.model.httpVerificationScheme
}

func (c *ServiceConfig) SenderEmail() string {
	return c.configEnv.SenderEmail
}

func (c *ServiceConfig) SenderPassword() string {
	return c.configEnv.SenderPassword
}

func (c *ServiceConfig) HostSMTP() string {
	return c.configEnv.HostSMTP
}

func (c *ServiceConfig) PortSMTP() string {
	return c.configEnv.PortSMTP
}

func (c *ServiceConfig) LoadFromMap(config map[string]json.RawMessage) error {
	c.model.mu.Lock()
	defer c.model.mu.Unlock()

	crt := reflect.TypeOf(&c.model).Elem()
	crv := reflect.ValueOf(&c.model).Elem()

	for idx := 0; idx < crt.NumField(); idx++ {

		crtf := crt.Field(idx)
		nameTag := crtf.Tag.Get("setting_name")
		tagDefault, hasTagDefault := crtf.Tag.Lookup("default")
		if nameTag == "" {
			continue
		}

		crvf := crv.Field(idx)
		value, ok := config[nameTag]
		if !ok {
			if !hasTagDefault {
				continue
			}
			if !crvf.IsZero() {
				continue
			}
			value = json.RawMessage(tagDefault)
			config[nameTag] = value
		}

		space := reflect.New(crvf.Type())
		intrfc := space.Interface()
		if err := json.Unmarshal(value, &intrfc); err != nil {
			log.Printf("failed to unmarshal value %s\n", nameTag)
			continue
		}
		crvf.Set(space.Elem())
	}

	return nil
}

func (c *ServiceConfig) ToMap() map[string]json.RawMessage {
	c.model.mu.RLock()
	defer c.model.mu.RUnlock()

	crt := reflect.TypeOf(&c.model).Elem()
	crv := reflect.ValueOf(&c.model).Elem()

	res := make(map[string]json.RawMessage, crt.NumField()-1)

	for idx := 0; idx < crt.NumField(); idx++ {
		crtf := crt.Field(idx)
		tag, ok := crtf.Tag.Lookup("setting_name")
		if !ok {
			continue
		}
		bts, _ := json.Marshal(crv.Field(idx).Interface())
		res[tag] = bts
	}
	return res
}

func (c *ServiceConfig) Validate(cfg map[string]json.RawMessage) error {
	c.model.mu.Lock()
	defer c.model.mu.Unlock()

	crt := reflect.TypeOf(&c.model).Elem()
	crv := reflect.ValueOf(&c.model).Elem()

	for idx := 0; idx < crt.NumField(); idx++ {

		crtf := crt.Field(idx)
		tag, ok := crtf.Tag.Lookup("setting_name")
		if !ok {
			continue
		}

		value, cfgHasTag := cfg[tag]
		if !cfgHasTag {
			continue
		}

		crvf := crv.Field(idx)

		space := reflect.New(crvf.Type())
		intrfc := space.Interface()
		if err := json.Unmarshal(value, &intrfc); err != nil {
			return errors.Errorf("failed to unmarshal value %s\n", tag)
		}
	}
	return nil
}
