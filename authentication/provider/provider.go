package provider

import (
	"slices"
)

type Provider struct {
	slug    string
	display string
}

func (p Provider) String() string {
	return p.slug
}

func (p Provider) Pretty() string {
	return p.display
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

func FromSlug(p string) (Provider, bool) {

	for _, single := range all {
		if single.slug == p {
			return single, true
		}
	}

	return Unknown, false

}

func List() []Provider {
	return all
}

func Slugs() []string {

	list := []string{}

	for _, single := range List() {
		list = append(list, single.slug)
	}

	return list
}

func Public() []string {

	list := []string{}

	for _, single := range List() {
		list = append(list, single.display)
	}

	return list
}

func Validate(p Provider) bool {
	return slices.Contains(List(), p)
}
