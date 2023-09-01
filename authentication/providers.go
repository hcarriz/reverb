package authentication

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/patreon"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/strava"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/tiktok"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitterv2"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"github.com/markbates/goth/providers/zoom"
)

type Provider struct {
	slug    string
	display string
}

func (p Provider) String() string {
	return p.slug
}

var (
	Unknown      = Provider{"", ""}
	Amazon       = Provider{"amazon", "Amazon"}
	Apple        = Provider{"apple", "Apple"}
	Auth0        = Provider{"auth0", "Auth0"}
	Azure        = Provider{"azure_ad", "Azure AD"}
	Battlenet    = Provider{"battlenet", "Battlenet"}
	Bitbucket    = Provider{"bitbucket", "Bitbucket"}
	Box          = Provider{"box", "Box"}
	Dailymotion  = Provider{"dailymotion", "Dailymotion"}
	Deezer       = Provider{"deezer", "Deezer"}
	DigitalOcean = Provider{"digital_ocean", "Digital Ocean"}
	Discord      = Provider{"discord", "Discord"}
	Dropbox      = Provider{"dropbox", "Dropbox"}
	Eve          = Provider{"eve_online", "Eve Online"}
	Facebook     = Provider{"facebook", "Facebook"}
	Fitbit       = Provider{"fitbit", "Fitbit"}
	Gitea        = Provider{"gitea", "Gitea"}
	Github       = Provider{"github", "Github"}
	Gitlab       = Provider{"gitlab", "Gitlab"}
	Google       = Provider{"google", "Google"}
	GooglePlus   = Provider{"google_plus", "Google Plus"}
	Heroku       = Provider{"heroku", "Heroku"}
	Instagram    = Provider{"instagram", "Instagram"}
	Intercom     = Provider{"intercom", "Intercom"}
	Kakao        = Provider{"kakao", "Kakao"}
	LastFM       = Provider{"last_fm", "Last FM"}
	LINE         = Provider{"line", "LINE"}
	Linkedin     = Provider{"linkedin", "Linkedin"}
	Mastodon     = Provider{"mastodon", "Mastodon"}
	Meetup       = Provider{"meetup", "Meetup.com"}
	Microsoft    = Provider{"microsoft_online", "Microsoft Online"}
	Naver        = Provider{"naver", "Naver"}
	NextCloud    = Provider{"nextCloud", "NextCloud"}
	Okta         = Provider{"okta", "Okta"}
	Onedrive     = Provider{"onedrive", "Onedrive"}
	OpenID       = Provider{"openid_connect", "OpenID Connect"}
	Patreon      = Provider{"patreon", "Patreon"}
	Paypal       = Provider{"paypal", "Paypal"}
	Salesforce   = Provider{"salesforce", "Salesforce"}
	SeaTalk      = Provider{"seatalk", "SeaTalk"}
	Shopify      = Provider{"shopify", "Shopify"}
	Slack        = Provider{"slack", "Slack"}
	SoundCloud   = Provider{"soundcloud", "SoundCloud"}
	Spotify      = Provider{"spotify", "Spotify"}
	Steam        = Provider{"steam", "Steam"}
	Strava       = Provider{"strava", "Strava"}
	Stripe       = Provider{"stripe", "Stripe"}
	TikTok       = Provider{"tikTok", "TikTok"}
	Twitch       = Provider{"twitch", "Twitch"}
	Twitter      = Provider{"twitter", "Twitter"}
	Typetalk     = Provider{"typetalk", "Typetalk"}
	Uber         = Provider{"uber", "Uber"}
	VK           = Provider{"vk", "VK"}
	WeCom        = Provider{"wecom", "WeCom"}
	Wepay        = Provider{"wepay", "Wepay"}
	Xero         = Provider{"xero", "Xero"}
	Yahoo        = Provider{"yahoo", "Yahoo"}
	Yammer       = Provider{"yammer", "Yammer"}
	Yandex       = Provider{"yandex", "Yandex"}
	Zoom         = Provider{"zoom", "Zoom"}
)

var all = []Provider{
	Amazon,
	Apple,
	Auth0,
	Azure,
	Battlenet,
	Bitbucket,
	Box,
	Dailymotion,
	Deezer,
	DigitalOcean,
	Discord,
	Dropbox,
	Eve,
	Facebook,
	Fitbit,
	Gitea,
	Github,
	Gitlab,
	Google,
	GooglePlus,
	Heroku,
	Instagram,
	Intercom,
	Kakao,
	LastFM,
	LINE,
	Linkedin,
	Mastodon,
	Meetup,
	Microsoft,
	Naver,
	NextCloud,
	Okta,
	Onedrive,
	OpenID,
	Patreon,
	Paypal,
	Salesforce,
	SeaTalk,
	Shopify,
	Slack,
	SoundCloud,
	Spotify,
	Steam,
	Strava,
	Stripe,
	TikTok,
	Twitch,
	Twitter,
	Typetalk,
	Uber,
	VK,
	WeCom,
	Wepay,
	Xero,
	Yahoo,
	Yammer,
	Yandex,
	Zoom,
}

func Slugs() []string {

	list := []string{}

	for _, single := range all {
		list = append(list, single.slug)
	}

	return list
}

func Public() []string {

	list := []string{}

	for _, single := range all {
		list = append(list, single.display)
	}

	return list
}

func Validate(p Provider) bool {
	return slices.Contains(all, p)
}

// WithProvider adds a provider. Source is required for Okta, Nextcloud, and OpenID Providers.
func WithProvider(p Provider, key, secret, callbackDomain, source string) Option {
	return option(func(a *auth) error {

		if !Validate(p) {
			return ErrInvalidProvider
		}

		if slices.Contains([]Provider{Okta, NextCloud, OpenID}, p) {
			if _, err := url.Parse(source); err != nil {
				return err
			}
		}

		u, err := url.Parse(callbackDomain)
		if err != nil {
			return err
		}

		u.Path = fmt.Sprintf("/auth/callback/%s", p.slug)

		d := cloneURL(u)
		d.Path = ""

		switch p {
		case Amazon:
			goth.UseProviders(amazon.New(key, secret, u.String()))
		case Apple:
			goth.UseProviders(apple.New(key, secret, u.String(), nil, apple.ScopeEmail, apple.ScopeName))
		case Auth0:
			goth.UseProviders(auth0.New(key, secret, u.String(), d.String()))
		case Azure:
			goth.UseProviders(azuread.New(key, secret, u.String(), nil))
		case Battlenet:
			goth.UseProviders(battlenet.New(key, secret, u.String()))
		case Bitbucket:
			goth.UseProviders(bitbucket.New(key, secret, u.String()))
		case Box:
			goth.UseProviders(box.New(key, secret, u.String()))
		case Dailymotion:
			goth.UseProviders(dailymotion.New(key, secret, u.String(), "email"))
		case Deezer:
			goth.UseProviders(deezer.New(key, secret, u.String(), "email"))
		case DigitalOcean:
			goth.UseProviders(digitalocean.New(key, secret, u.String(), "read"))
		case Discord:
			goth.UseProviders(discord.New(key, secret, u.String(), discord.ScopeIdentify, discord.ScopeEmail))
		case Dropbox:
			goth.UseProviders(dropbox.New(key, secret, u.String()))
		case Eve:
			goth.UseProviders(eveonline.New(key, secret, u.String()))
		case Facebook:
			goth.UseProviders(facebook.New(key, secret, u.String()))
		case Fitbit:
			goth.UseProviders(fitbit.New(key, secret, u.String()))
		case Gitea:
			goth.UseProviders(gitea.New(key, secret, u.String()))
		case Github:
			goth.UseProviders(github.New(key, secret, u.String(), "read:user", "user:email"))
		case Gitlab:
			goth.UseProviders(gitlab.New(key, secret, u.String()))
		case Google:
			goth.UseProviders(google.New(key, secret, u.String()))
		case GooglePlus:
			goth.UseProviders(gplus.New(key, secret, u.String()))
		case Heroku:
			goth.UseProviders(heroku.New(key, secret, u.String()))
		case Instagram:
			goth.UseProviders(instagram.New(key, secret, u.String()))
		case Intercom:
			goth.UseProviders(intercom.New(key, secret, u.String()))
		case Kakao:
			goth.UseProviders(kakao.New(key, secret, u.String()))
		case LastFM:
			goth.UseProviders(lastfm.New(key, secret, u.String()))
		case LINE:
			goth.UseProviders(line.New(key, secret, u.String(), "profile", "openid", "email"))
		case Linkedin:
			goth.UseProviders(linkedin.New(key, secret, u.String()))
		case Mastodon:
			goth.UseProviders(mastodon.New(key, secret, u.String(), "read:accounts"))
		case Meetup:
			goth.UseProviders(meetup.New(key, secret, u.String()))
		case Microsoft:
			goth.UseProviders(microsoftonline.New(key, secret, u.String()))
		case Naver:
			goth.UseProviders(naver.New(key, secret, u.String()))
		case NextCloud:
			goth.UseProviders(nextcloud.NewCustomisedDNS(key, secret, u.String(), source))
		case Okta:
			goth.UseProviders(okta.New(key, secret, source, u.String()))
		case Onedrive:
			goth.UseProviders(onedrive.New(key, secret, u.String()))
		case OpenID:

			oic, err := openidConnect.New(key, secret, u.String(), source)
			if err != nil {
				return err
			}

			goth.UseProviders(oic)

		case Patreon:
			goth.UseProviders(patreon.New(key, secret, u.String()))
		case Paypal:
			goth.UseProviders(paypal.New(key, secret, u.String()))
		case Salesforce:
			goth.UseProviders(salesforce.New(key, secret, u.String()))
		case SeaTalk:
			goth.UseProviders(seatalk.New(key, secret, u.String()))
		case Shopify:
			goth.UseProviders(shopify.New(key, secret, u.String()))
		case Slack:
			goth.UseProviders(slack.New(key, secret, u.String()))
		case SoundCloud:
			goth.UseProviders(soundcloud.New(key, secret, u.String()))
		case Spotify:
			goth.UseProviders(spotify.New(key, secret, u.String()))
		case Steam:
			goth.UseProviders(steam.New(key, u.String()))
		case Strava:
			goth.UseProviders(strava.New(key, secret, u.String()))
		case Stripe:
			goth.UseProviders(stripe.New(key, secret, u.String()))
		case TikTok:
			goth.UseProviders(tiktok.New(key, secret, u.String()))
		case Twitch:
			goth.UseProviders(twitch.New(key, secret, u.String()))
		case Twitter:
			goth.UseProviders(twitterv2.New(key, secret, u.String()))
		case Typetalk:
			goth.UseProviders(typetalk.New(key, secret, u.String(), "my"))
		case Uber:
			goth.UseProviders(uber.New(key, secret, u.String()))
		case VK:
			goth.UseProviders(vk.New(key, secret, u.String()))
		case Wepay:
			goth.UseProviders(wepay.New(key, secret, u.String(), "view_user"))
		case Xero:
			goth.UseProviders(xero.New(key, secret, u.String()))
		case Yahoo:
			goth.UseProviders(yahoo.New(key, secret, u.String()))
		case Yammer:
			goth.UseProviders(yammer.New(key, secret, u.String()))
		case Yandex:
			goth.UseProviders(yandex.New(key, secret, u.String()))
		case Zoom:
			goth.UseProviders(zoom.New(key, secret, u.String(), "read:user"))
		default:
			return ErrInvalidProvider
		}

		a.providers = append(a.providers, p)

		return nil
	})
}
