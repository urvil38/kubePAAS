package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

func GetStorageClient() (storageClient *storage.Client, err error) {

	config := jwt.Config{
		PrivateKeyID: "ff5a14b807fa7c77e1110b3ee9f373ce6dfea0af",
		Email:        "storagekubepaas@kubepaas.iam.gserviceaccount.com",
		PrivateKey:   []byte("-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCXZMEdJ0AYbW/Q\nLmoQfPzfN8xnJDRbsdYvB1w/6yuQ2Kah+cLEB0rFpFccD7HwSWjwVQ0UcHTPoFMs\nXhmbT14SY3i5eejTIutZZYpzQlC6uGmLoVwmzhrPZmy4wI9OHH0GEBBAyFUQnHQ+\nk7jO7oFzBghxZIvPC5BM2uZQxYhXKQz2vlM3gSRT2STwDlsn9sMiotoJ0pSmAg5N\nGa/OrN9bytgXPkdsFpAG+94L/gxD1UPGfnh2HB8kCPNV5zK5Yy5UukOkoq+0g059\nVAfRRJXYWNbaCm1gUBMk6aICa0fB6aamSNAN7rFFx7IeLBG0vF+d/gGpUyUqiuD6\n9EjqUXflAgMBAAECggEAAOAhScTfVwTS/7Y1ANFoOPY+pV4NO3aE0ZLOUsROZTEL\njaY/HRkZspjntA7XLZePFsy3HaYk1sqLkJceuMo1tg+DNDdjRE1QZRz0NwRsKRhF\n6/vL56GLgCWMfWkHqyD5DB7tqSI/c7Aj7/S0veWdNAgV5mn5cQIVHIyrhk4OIsrk\ns81Qq3xx8UXJdR1J+2bZJtMm1xvN0I0sVRE2VPVXe/mhU3ogxZZGcAk2UzSP/+o/\nIV8uwaoa6rux9LWJC/2Di5JMwGC69un2c7Q7Tm5EjVuuTMgTyp8J3DA2cj1kTiji\nzRuiGMtyFrSnKeU5pKdac7jxBxyt8cIyFoLpX2B9IQKBgQDN45smnlTBsNewOTPb\nDtzkcck756SL9LhHBV/CHa4vAjplNzWq70urP/r7lTqwsOUmjlJsUpQWn6KZ+wI1\n2F7wdJG5+GZpLni99dDMCWCH7880X5935moxZ6UQnaKT4c41Qf2f17LFV/EoF91F\nxpc0ejZ1mOFilhJhOtSiWfT3nQKBgQC8Pa4AJvFUpIoU3YatkBjSE/AiWSHtbxnm\nMGSM/DEMrJGYoOiyS3O+EN2ESu8fz5Aq/IUfJ613ZvEHWeMgUsp3O3NLv3+XJNLt\n99pcbleR4EZRyTEoVuOHv+unYNcx0HFAwz2Kb0gonk7cD4Bqvt1/hXSiUXyHhNcq\n1xbysjpi6QKBgG5oR6MF9N2JP6C4jB5Ech/vBMKjDZIfwVIUh61IzUdYIoh0essX\nV7SVsrmG5NorgjaSy5BbGB6prEB2YlENnpvDZwIbOo/c49K6JyXDQYikCLFFNfbO\nENQ9iD7IyY4T4Miegqtct/krl56wbXAldqAliV62hOahI2oQakZFhx/hAoGAFVr+\nAYBpgovEKofTPp+JYVPnu03XXoNrMcUtsxztR53QSrt+irOqptZs+xQxOq+mkGnj\nhFxQ/qnMEGRvMvyRgaNZ+i74f6Iq19p1iGTwRFloQOENVaE94OmyB2QiJcGbB5je\nw3TKf+kt0yNjEzkiEdSHHd8WvJ0id/a84L29h/kCgYAhnMSM9D04sLbi8/a+I6lB\naSE4DrNeE8Yp+BAc6ucQrcxex471breLqIe80RLBX/EZQm4eHrBLP/EwgtNbAFVS\nHAuW29ZD18AGuZf070kZ/pj2Cx3tFKhe8e91TIJJpqAHNYZrGhhDp51omx+uoPzX\nt6h7dXtjj9v28KbLBgAwQA==\n-----END PRIVATE KEY-----\n"),
		Scopes:       []string{"https://www.googleapis.com/auth/devstorage.read_write"},
		TokenURL:     google.JWTTokenURL,
	}

	httpClient := config.Client(oauth2.NoContext)
	storageClient, err = storage.NewClient(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}