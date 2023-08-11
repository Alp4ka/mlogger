package mlogger

import (
	"context"
	"fmt"
	"github.com/Alp4ka/mlogger/contactpoints/matrix"
	"github.com/Alp4ka/mlogger/field"
	"github.com/Alp4ka/mlogger/jsonsecurity"
	"github.com/Alp4ka/mlogger/misc"
	"github.com/Alp4ka/mlogger/templates"
	"sync"
	"testing"
	"time"
)

const t = `
My test template! 

**Time:**
*{{ .LogTime.Format "Jan 02, 2006 15:04:05 UTC" }}*

**Level:**
*{{ .LogLevel }}*

**Origin:**
*{{ .LogSource }}*
{{ if len .LogContextFields }}
**Context Fields:**
{{- end }}
{{ range .LogContextFields }}
{{ .Key }}: {{ .Value }}
{{ end }}
{{ if len .LogFields }}
**Fields:**
{{- end }}
{{ range .LogFields }}
{{ .Key }}: {{ .Value }}
{{ end }}
**Message:**
*{{ .LogMessage }}*`

func Test_Main(test *testing.T) {
	m := matrix.NewContactPoint(matrix.Config{})
	cfg := Config{Level: misc.LevelInfo,
		Template: templates.Config{Pattern: t, Use: true}}

	f1 := field.Bool("test_bool", true)
	ctx := field.WithContextFields(context.Background(), f1)
	logger, err := NewProduction(ctx, cfg, m)
	if err != nil {
		panic(err.Error())
	}

	ReplaceGlobals(logger)
	L().Info(
		"test message",
		field.Int("test_int", 123),
		field.String("test_string", "hello world!"),
		field.Error(fmt.Errorf("test_error")))

	time.Sleep(time.Second)
}

func Test_JSONField(test *testing.T) {
	cfg := Config{Level: misc.LevelInfo}
	logger, err := NewProduction(context.Background(), cfg)
	if err != nil {
		test.Errorf("Failed to init logger: %v", err)
	}

	logger.Info("json test msg", field.JSONEscape("my_json", []byte("{\"int_field\":123, \"string_field\":\"hello\"}")))
}

func Test_JSONSecureField(test *testing.T) {
	cfg := Config{
		Level: misc.LevelInfo,
		JSONSecurity: jsonsecurity.Config{
			MaxDepth: 10,
			Triggers: map[string]jsonsecurity.TriggerOpts{
				"email":          {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelEmail},
				"e-mail":         {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelEmail},
				"otp":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"otp-code":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"otpcode":        {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"security":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"cvv":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCVV},
				"cvc":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCVV},
				"cardholder":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"cardholdername": {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"ifsccode":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPhoneNumber},
				"phonenumber":    {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPhoneNumber},
				"accountnumber":  {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"iban":           {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"clientname":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"card":           {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"upiid":          {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"cardnumber":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"recipient":      {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"credentials":    {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"password":       {CaseSensitive: true, ShouldAppear: false, MaskMethod: jsonsecurity.MaskerLabelPassword},
			},
		},
	}
	logger, err := NewProduction(context.Background(), cfg)
	if err != nil {
		test.Errorf("Failed to init logger: %v", err)
	}

	logger.Info(
		"json test msg",
		field.JSONEscapeSecure(
			"my_secure_json",
			[]byte(`{
  "email": "example123@example.com",
  "UUID": "209e9061-72a7-48c1-a020-92ba6dcb6c71",
  "clientOrderID": "12345qwerty",
  "amount": 30000,
  "password": "my_password",
  "expireAt": 3600,
  "comment": "some comment",
  "currencyID": 1,
  "callbackURL": "https://google.com",
  "recipient": "123456789012334",
  "clientInfo": {
	"password": "my_password_here",
    "client": "710",
    "recipient_system": "string"
  }
}`,
			),
		),
	)
}

func Test_JSONSecureFieldParallel(test *testing.T) {
	cfg := Config{
		Level: misc.LevelInfo,
		JSONSecurity: jsonsecurity.Config{
			MaxDepth: 10,
			Triggers: map[string]jsonsecurity.TriggerOpts{
				"email":          {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelEmail},
				"e-mail":         {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelEmail},
				"otp":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"otp-code":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"otpcode":        {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"security":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
				"cvv":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCVV},
				"cvc":            {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCVV},
				"cardholder":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"cardholdername": {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"ifsccode":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPhoneNumber},
				"phonenumber":    {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPhoneNumber},
				"accountnumber":  {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"iban":           {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"clientname":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelName},
				"card":           {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"upiid":          {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"cardnumber":     {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"recipient":      {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"credentials":    {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelCardNumber},
				"password":       {CaseSensitive: false, ShouldAppear: true, MaskMethod: jsonsecurity.MaskerLabelPassword},
			},
		},
	}
	logger, err := NewProduction(context.Background(), cfg)
	if err != nil {
		test.Errorf("Failed to init logger: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			logger.Info(
				"json test msg",
				field.JSONEscapeSecure(
					"my_secure_json",
					[]byte(`{
  "email": "example123@example.com",
  "UUID": "209e9061-72a7-48c1-a020-92ba6dcb6c71",
  "clientOrderID": "12345qwerty",
  "amount": 30000,
  "password": "my_password",
  "expireAt": 3600,
  "comment": "some comment",
  "currencyID": 1,
  "callbackURL": "https://google.com",
  "recipient": "123456789012334",
  "clientInfo": {
	"password": ["my_password_here", 123],
    "client": "710",
    "recipient_system": "string"
  }
}`,
					),
				),
			)

			wg.Done()
		}()
	}
	wg.Wait()

	logger.Info("got it")
}

func Test_Parallel(test *testing.T) {
	cfg := Config{
		Level: misc.LevelInfo,
	}
	logger, err := NewProduction(context.Background(), cfg)
	if err != nil {
		test.Errorf("Failed to init logger: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			logger.Info(
				"json test msg",
				field.JSONEscape(
					"my_secure_json",
					[]byte(`{
  "email": "example123@example.com",
  "UUID": "209e9061-72a7-48c1-a020-92ba6dcb6c71",
  "clientOrderID": "12345qwerty",
  "amount": 30000,
  "password": "my_password",
  "expireAt": 3600,
  "comment": "some comment",
  "currencyID": 1,
  "callbackURL": "https://google.com",
  "recipient": "123456789012334",
  "clientInfo": {
	"password": ["my_password_here", 123],
    "client": "710",
    "recipient_system": "string"
  }
}`,
					),
				),
			)

			wg.Done()
		}()
	}
	wg.Wait()

	logger.Info("got it")
}
