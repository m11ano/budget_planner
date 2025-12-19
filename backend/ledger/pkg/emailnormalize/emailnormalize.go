package emailnormalize

import (
	"errors"
	"net/mail"
	"strings"
)

type Rules uint

const (
	DashAddressing Rules = 1 << iota
	PlusAddressing
	LocalPartAsHostname
	StripPeriods
)

type Provider struct {
	Name      string
	Domains   []string
	Canonical string
	Flags     Rules
}

type Result struct {
	Address           string
	NormalizedAddress string
	Provider          string
}

var Providers = []Provider{
	{
		Name:      "Fastmail",
		Domains:   []string{"fastmail.com", "messagingengine.com", "fastmail.fm"},
		Canonical: "fastmail.com",
		Flags:     PlusAddressing | LocalPartAsHostname,
	},
	{
		Name:      "Apple",
		Domains:   []string{"icloud.com", "me.com", "mac.com"},
		Canonical: "icloud.com",
		Flags:     PlusAddressing,
	},
	{
		Name: "Yahoo",
		Domains: []string{
			"yahoo.com.ar", "yahoo.com.au", "yahoo.at", "yahoo.be", "yahoo.com.br",
			"ca.yahoo.com", "qc.yahoo.com", "yahoo.com.co", "yahoo.com.hr", "yahoo.cz",
			"yahoo.dk", "yahoo.fi", "yahoo.fr", "yahoo.de", "yahoo.gr",
			"yahoo.com.hk", "yahoo.hu", "yahoo.co.in", "yahoo.in", "yahoo.co.id",
			"yahoo.ie", "yahoo.co.il", "yahoo.it", "yahoo.co.jp", "yahoo.com.my",
			"yahoo.com.mx", "yahoo.ae", "yahoo.nl", "yahoo.co.nz", "yahoo.no",
			"yahoo.com.ph", "yahoo.pl", "yahoo.pt", "yahoo.ro", "yahoo.ru",
			"yahoo.com.sg", "yahoo.co.za", "yahoo.es", "yahoo.se", "yahoo.ch/fr",
			"yahoo.ch/de", "yahoo.com.tw", "yahoo.co.th", "yahoo.com.tr", "yahoo.co.uk",
			"yahoo.com", "yahoo.com.vn", "ymail.com", "yahoodns.net",
		},
		Canonical: "yahoo.com",
		Flags:     DashAddressing,
	},
	{
		Name:      "Google",
		Domains:   []string{"gmail.com", "googlemail.com", "google.com"},
		Canonical: "gmail.com",
		Flags:     PlusAddressing | StripPeriods,
	},
	{
		Name:      "Rambler",
		Domains:   []string{"rambler.ru", "lenta.ru", "autorambler.ru", "myrambler.ru", "ro.ru"},
		Canonical: "rambler.ru",
		Flags:     0,
	},
	{
		Name: "Microsoft",
		Domains: []string{
			"hotmail.com", "hotmail.at", "hotmail.be", "hotmail.ca", "hotmail.cl",
			"hotmail.co.il", "hotmail.co.nz", "hotmail.co.th", "hotmail.co.uk",
			"hotmail.com.ar", "hotmail.com.au", "hotmail.com.br", "hotmail.com.gr",
			"hotmail.com.mx", "hotmail.com.pe", "hotmail.com.tr", "hotmail.com.vn",
			"hotmail.cz", "hotmail.de", "hotmail.dk", "hotmail.es", "hotmail.fr",
			"hotmail.hu", "hotmail.id", "hotmail.ie", "hotmail.in", "hotmail.it",
			"hotmail.jp", "hotmail.kr", "hotmail.lv", "hotmail.my", "hotmail.ph",
			"hotmail.pt", "hotmail.sa", "hotmail.sg", "hotmail.sk", "live.com",
			"live.be", "live.co.uk", "live.com.ar", "live.com.mx", "live.de",
			"live.es", "live.eu", "live.fr", "live.it", "live.nl", "msn.com",
			"outlook.com", "outlook.at", "outlook.be", "outlook.cl", "outlook.co.il",
			"outlook.co.nz", "outlook.co.th", "outlook.com.ar", "outlook.com.au",
			"outlook.com.br", "outlook.com.gr", "outlook.com.pe", "outlook.com.tr",
			"outlook.com.vn", "outlook.cz", "outlook.de", "outlook.dk", "outlook.es",
			"outlook.fr", "outlook.hu", "outlook.id", "outlook.ie", "outlook.in",
			"outlook.it", "outlook.jp", "outlook.kr", "outlook.lv", "outlook.my",
			"outlook.ph", "outlook.pt", "outlook.sa", "outlook.sg", "outlook.sk",
			"passport.com",
		},
		Canonical: "outlook.com",
		Flags:     PlusAddressing,
	},
	{
		Name: "Yandex",
		Domains: []string{
			"narod.ru", "yandex.ru", "yandex.org", "yandex.net", "yandex.net.ru",
			"yandex.com.ru", "yandex.ua", "yandex.com.ua", "yandex.by", "yandex.eu",
			"yandex.ee", "yandex.lt", "yandex.lv", "yandex.md", "yandex.uz",
			"yandex.mx", "yandex.do", "yandex.tm", "yandex.de", "yandex.ie",
			"yandex.in", "yandex.qa", "yandex.so", "yandex.nu", "yandex.tj",
			"yandex.dk", "yandex.es", "yandex.pt", "yandex.kz", "yandex.pl",
			"yandex.lu", "yandex.it", "yandex.az", "yandex.ro", "yandex.rs",
			"yandex.sk", "yandex.no", "ya.ru", "yandex.com", "yandex.asia",
			"yandex.mobi",
		},
		Canonical: "yandex.ru",
		Flags:     PlusAddressing,
	},
	{
		Name:      "ProtonMail",
		Domains:   []string{"protonmail.ch", "protonmail.com", "proton.me", "pm.me"},
		Canonical: "protonmail.ch",
		Flags:     PlusAddressing,
	},
}

func Normalize(email string) (*Result, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}
	raw := strings.ToLower(addr.Address)
	parts := strings.SplitN(raw, "@", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid email address")
	}
	local, domain := parts[0], parts[1]

	provider := lookupProvider(domain)
	if provider != nil {
		if provider.Flags&LocalPartAsHostname != 0 {
			if i := strings.Index(domain, "."); i != -1 {
				local = domain[:i]
			}
		}

		domain = provider.Canonical

		if provider.Flags&StripPeriods != 0 {
			local = strings.ReplaceAll(local, ".", "")
		}

		if provider.Flags&PlusAddressing != 0 {
			if i := strings.Index(local, "+"); i != -1 {
				local = local[:i]
			}
		}

		if provider.Flags&DashAddressing != 0 {
			if i := strings.Index(local, "-"); i != -1 {
				local = local[:i]
			}
		}
	}

	return &Result{
		Address:           addr.Address,
		NormalizedAddress: local + "@" + domain,
		Provider: func() string {
			if provider != nil {
				return provider.Name
			}
			return ""
		}(),
	}, nil
}

func lookupProvider(domain string) *Provider {
	for i := range Providers {
		p := &Providers[i]
		for _, d := range p.Domains {
			if domain == d {
				return p
			}
		}
	}
	return nil
}
