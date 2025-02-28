package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/header"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("PHOTOPRISM_TEST", "true")
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	c := TestConfig()
	defer c.CloseDb()

	code := m.Run()

	os.Exit(code)
}

func TestNewConfig(t *testing.T) {
	ctx := CliTestContext()

	assert.True(t, ctx.IsSet("assets-path"))
	assert.False(t, ctx.Bool("debug"))

	c := NewConfig(ctx)

	assert.IsType(t, new(Config), c)

	assert.Equal(t, fs.Abs("../../assets"), c.AssetsPath())
	assert.False(t, c.Prod())
	assert.False(t, c.Debug())
	assert.False(t, c.ReadOnly())
}

func TestConfig_Prod(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.False(t, c.Prod())
	assert.False(t, c.Debug())
	assert.False(t, c.Trace())
	c.options.Prod = true
	c.options.Debug = true
	assert.True(t, c.Prod())
	assert.False(t, c.Debug())
	assert.False(t, c.Trace())
	c.options.Prod = false
	assert.True(t, c.Debug())
	assert.False(t, c.Trace())
	c.options.Debug = false
	assert.False(t, c.Debug())
	assert.False(t, c.Debug())
	assert.False(t, c.Trace())
}

func TestConfig_Name(t *testing.T) {
	c := NewConfig(CliTestContext())

	name := c.Name()
	assert.Equal(t, "PhotoPrism", name)
}

func TestConfig_About(t *testing.T) {
	c := NewConfig(CliTestContext())

	name := c.About()
	assert.Equal(t, "PhotoPrism®", name)
}

func TestConfig_Edition(t *testing.T) {
	c := NewConfig(CliTestContext())

	name := c.Edition()
	assert.NotEmpty(t, name)
}

func TestConfig_Version(t *testing.T) {
	c := NewConfig(CliTestContext())

	version := c.Version()
	assert.Equal(t, "0.0.0", version)
}

func TestConfig_VersionChecksum(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, uint32(0x2e5b4b86), c.VersionChecksum())
}

func TestConfig_TensorFlowVersion(t *testing.T) {
	c := NewConfig(CliTestContext())

	version := c.TensorFlowVersion()
	assert.IsType(t, "1.15.0", version)
}

func TestConfig_TensorFlowDisabled(t *testing.T) {
	c := NewConfig(CliTestContext())

	version := c.DisableTensorFlow()
	assert.Equal(t, false, version)
}

func TestConfig_Copyright(t *testing.T) {
	c := NewConfig(CliTestContext())

	copyright := c.Copyright()
	assert.Equal(t, "", copyright)
}

func TestConfig_OptionsYaml(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		assert.Contains(t, c.OptionsYaml(), "options.yml")
	})

	t.Run("ChangePath", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		assert.Contains(t, c.OptionsYaml(), "options.yml")
		c.options.ConfigPath = "/go/src/github.com/photoprism/photoprism/internal/config/testdata/"
		assert.Equal(t, "/go/src/github.com/photoprism/photoprism/internal/config/testdata/options.yml", c.OptionsYaml())
	})
}

func TestConfig_PIDFilename(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Contains(t, c.PIDFilename(), "/storage/testdata/photoprism.pid")
}

func TestConfig_LogFilename(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Contains(t, c.LogFilename(), "/storage/testdata/photoprism.log")
}

func TestConfig_DetachServer(t *testing.T) {
	c := NewConfig(CliTestContext())

	detachServer := c.DetachServer()
	assert.Equal(t, false, detachServer)
}

func TestConfig_OriginalsPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	result := c.OriginalsPath()
	assert.True(t, strings.HasPrefix(result, "/"))
	assert.True(t, strings.HasSuffix(result, "/storage/testdata/originals"))
}

func TestConfig_ImportPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	result := c.ImportPath()
	assert.True(t, strings.HasPrefix(result, "/"))
	assert.True(t, strings.HasSuffix(result, "/storage/testdata/import"))
}

func TestConfig_CachePath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.HasSuffix(c.CachePath(), "storage/testdata/cache"))
}

func TestConfig_MediaCachePath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.HasPrefix(c.MediaCachePath(), "/"))
	assert.True(t, strings.HasSuffix(c.MediaCachePath(), "storage/testdata/cache/media"))
}

func TestConfig_MediaFileCachePath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, c.MediaCachePath(), c.MediaFileCachePath(""))
	assert.Equal(t, filepath.Join(c.MediaCachePath(), "a"), c.MediaFileCachePath("a"))
	assert.Equal(t, filepath.Join(c.MediaCachePath(), "0", "b", "5"), c.MediaFileCachePath("0b57b50fe3f6d12bbbf5f1abda3ebcc8bb5ebcee"))
}

func TestConfig_ThumbCachePath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.HasPrefix(c.ThumbCachePath(), "/"))
	assert.True(t, strings.HasSuffix(c.ThumbCachePath(), "storage/testdata/cache/thumbnails"))
}

func TestConfig_AssetsPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.True(t, strings.HasSuffix(c.AssetsPath(), "/assets"))
}

func TestConfig_CustomAssetsPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.CustomAssetsPath())
}

func TestConfig_DetectNSFW(t *testing.T) {
	c := NewConfig(CliTestContext())

	result := c.DetectNSFW()
	assert.Equal(t, true, result)
}

func TestConfig_AdminUser(t *testing.T) {
	c := NewConfig(CliTestContext())

	c.options.AdminUser = "foo  "
	assert.Equal(t, "foo", c.AdminUser())
	c.options.AdminUser = "  Admin"
	assert.Equal(t, "admin", c.AdminUser())
}

func TestConfig_AdminPassword(t *testing.T) {
	c := NewConfig(CliTestContext())

	result := c.AdminPassword()
	assert.Equal(t, "photoprism", result)
}

func TestConfig_NSFWModelPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Contains(t, c.NSFWModelPath(), "/assets/nsfw")
}

func TestConfig_FaceNetModelPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Contains(t, c.FaceNetModelPath(), "/assets/facenet")
}

func TestConfig_ExamplesPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.ExamplesPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/examples", path)
}

func TestConfig_TensorFlowModelPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.TensorFlowModelPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/nasnet", path)
}

func TestConfig_TemplatesPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.TemplatesPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/templates", path)
}

func TestConfig_CustomTemplatesPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.CustomTemplatesPath()
	assert.Equal(t, "", path)
}

func TestConfig_TemplatesFiles(t *testing.T) {
	c := NewConfig(CliTestContext())

	files := c.TemplateFiles()

	t.Logf("TemplateFiles: %#v", files)
}

func TestConfig_StaticPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.StaticPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/static", path)
}

func TestConfig_StaticFile(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.StaticFile("video/404.mp4")
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/static/video/404.mp4", path)
}

func TestConfig_BuildPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.BuildPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/static/build", path)
}

func TestConfig_ImgPath(t *testing.T) {
	c := NewConfig(CliTestContext())

	path := c.ImgPath()
	assert.Equal(t, "/go/src/github.com/photoprism/photoprism/assets/static/img", path)
}

func TestConfig_IndexWorkers(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.GreaterOrEqual(t, c.IndexWorkers(), 1)
}

func TestConfig_IndexSchedule(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, DefaultIndexSchedule, c.IndexSchedule())
}

func TestConfig_WakeupInterval(t *testing.T) {
	c := NewConfig(CliTestContext())
	i := c.WakeupInterval()
	assert.Equal(t, "1h34m9s", c.WakeupInterval().String())
	c.options.WakeupInterval = 45
	assert.Equal(t, "1m0s", c.WakeupInterval().String())
	c.options.WakeupInterval = 0
	assert.Equal(t, "15m0s", c.WakeupInterval().String())
	c.options.WakeupInterval = 150
	assert.Equal(t, "2m30s", c.WakeupInterval().String())
	c.options.WakeupInterval = i
	assert.Equal(t, "1h34m9s", c.WakeupInterval().String())
}

func TestConfig_AutoIndex(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, -1*time.Second, c.AutoIndex())
}

func TestConfig_AutoImport(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, 2*time.Hour, c.AutoImport())
}

func TestConfig_GeoApi(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "places", c.GeoApi())
	c.options.DisablePlaces = true
	assert.Equal(t, "", c.GeoApi())
}

func TestConfig_OriginalsLimit(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, -1, c.OriginalsLimit())
	c.options.OriginalsLimit = 800
	assert.Equal(t, 800, c.OriginalsLimit())
}

func TestConfig_OriginalsByteLimit(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, int64(-1), c.OriginalsByteLimit())
	c.options.OriginalsLimit = 800
	assert.Equal(t, int64(838860800), c.OriginalsByteLimit())
}

func TestConfig_ResolutionLimit(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, DefaultResolutionLimit, c.ResolutionLimit())
	c.options.ResolutionLimit = 800
	assert.Equal(t, 800, c.ResolutionLimit())
	c.options.ResolutionLimit = 950
	assert.Equal(t, 900, c.ResolutionLimit())
	c.options.ResolutionLimit = 0
	assert.Equal(t, DefaultResolutionLimit, c.ResolutionLimit())
	c.options.ResolutionLimit = -1
	assert.Equal(t, -1, c.ResolutionLimit())
	c.options.Sponsor = false
	assert.Equal(t, -1, c.ResolutionLimit())
	c.options.Sponsor = true
	assert.Equal(t, -1, c.ResolutionLimit())
}

func TestConfig_BaseUri(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.BaseUri(""))
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, "", c.BaseUri(""))
	c.options.SiteUrl = "http://foo:2342/foo bar/"
	assert.Equal(t, "/foo%20bar", c.BaseUri(""))
	assert.Equal(t, "/foo%20bar/baz", c.BaseUri("/baz"))
}

func TestConfig_StaticUri(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "/static", c.StaticUri())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, "/static", c.StaticUri())
	c.options.SiteUrl = "http://foo:2342/foo/"
	assert.Equal(t, "/foo/static", c.StaticUri())
	c.options.CdnUrl = "http://foo:2342/bar"
	assert.Equal(t, "http://foo:2342/bar/foo"+StaticUri, c.StaticUri())
}

func TestConfig_ApiUri(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, ApiUri, c.ApiUri())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, ApiUri, c.ApiUri())
	c.options.SiteUrl = "http://foo:2342/foo/"
	assert.Equal(t, "/foo"+ApiUri, c.ApiUri())
}

func TestConfig_CdnUrl(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.options.SiteUrl)
	assert.Equal(t, "", c.CdnUrl(""))
	assert.True(t, c.NoCdn())
	assert.False(t, c.UseCdn())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, "/", c.CdnUrl("/"))
	c.options.CdnUrl = "http://foo:2342/foo/"
	assert.Equal(t, "http://foo:2342/foo", c.CdnUrl(""))
	assert.Equal(t, "http://foo:2342/foo/", c.CdnUrl("/"))
	assert.False(t, c.NoCdn())
	assert.True(t, c.UseCdn())
	c.options.SiteUrl = c.options.CdnUrl
	assert.Equal(t, "/", c.CdnUrl("/"))
	assert.Equal(t, "", c.CdnUrl(""))
	assert.True(t, c.NoCdn())
	assert.False(t, c.UseCdn())
	c.options.SiteUrl = ""
	assert.False(t, c.NoCdn())
	assert.True(t, c.UseCdn())
}

func TestConfig_CdnDomain(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.options.SiteUrl)
	assert.Equal(t, "", c.CdnDomain())
	c.options.CdnUrl = "http://superhost:2342/"
	assert.Equal(t, "superhost", c.CdnDomain())
	c.options.CdnUrl = "https://foo.bar.com:2342/foo/"
	assert.Equal(t, "foo.bar.com", c.CdnDomain())
	c.options.SiteUrl = c.options.CdnUrl
	assert.Equal(t, "", c.CdnDomain())
	c.options.SiteUrl = ""
	c.options.CdnUrl = "http:/invalid:2342/foo/"
	assert.Equal(t, "", c.CdnDomain())
	c.options.CdnUrl = ""
	assert.Equal(t, "", c.CdnDomain())
}

func TestConfig_CdnVideo(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.False(t, c.CdnVideo())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.False(t, c.CdnVideo())
	c.options.CdnUrl = "http://foo:2342/foo/"
	assert.False(t, c.CdnVideo())
	c.options.CdnVideo = true
	assert.True(t, c.CdnVideo())
	c.options.SiteUrl = c.options.CdnUrl
	assert.False(t, c.CdnVideo())
	c.options.SiteUrl = ""
	assert.True(t, c.CdnVideo())
	c.options.CdnVideo = false
	assert.False(t, c.CdnVideo())
	c.options.CdnUrl = ""
	assert.False(t, c.CdnVideo())
}

func TestConfig_CORSOrigin(t *testing.T) {
	c := NewConfig(CliTestContext())

	c.Options().CORSOrigin = ""
	assert.Equal(t, "", c.CORSOrigin())
	c.Options().CORSOrigin = "*"
	assert.Equal(t, "*", c.CORSOrigin())
	c.Options().CORSOrigin = "https://developer.mozilla.org"
	assert.Equal(t, "https://developer.mozilla.org", c.CORSOrigin())
	c.Options().CORSOrigin = ""
	assert.Equal(t, "", c.CORSOrigin())
}

func TestConfig_CORSHeaders(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.CORSHeaders())
	c.Options().CORSHeaders = header.DefaultAccessControlAllowHeaders
	assert.Equal(t, header.DefaultAccessControlAllowHeaders, c.CORSHeaders())
	c.Options().CORSHeaders = ""
	assert.Equal(t, "", c.CORSHeaders())
}

func TestConfig_CORSMethods(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "", c.CORSMethods())
	c.Options().CORSMethods = header.DefaultAccessControlAllowMethods
	assert.Equal(t, header.DefaultAccessControlAllowMethods, c.CORSMethods())
	c.Options().CORSMethods = ""
	assert.Equal(t, "", c.CORSMethods())
}

func TestConfig_ContentUri(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, ApiUri, c.ContentUri())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, ApiUri, c.ContentUri())
	c.options.CdnUrl = "http://foo:2342//"
	assert.Equal(t, "http://foo:2342"+ApiUri, c.ContentUri())
}

func TestConfig_VideoUri(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, ApiUri, c.VideoUri())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, ApiUri, c.VideoUri())
	c.options.CdnUrl = "http://foo:2342//"
	c.options.CdnVideo = true
	assert.Equal(t, "http://foo:2342"+ApiUri, c.VideoUri())
	c.options.CdnVideo = false
	assert.Equal(t, ApiUri, c.VideoUri())
}

func TestConfig_SiteUrl(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "http://localhost:2342/", c.SiteUrl())
	c.options.SiteUrl = "http://superhost:2342/"
	assert.Equal(t, "http://superhost:2342/", c.SiteUrl())
	c.options.SiteUrl = "http://superhost"
	assert.Equal(t, "http://superhost/", c.SiteUrl())
}

func TestConfig_SiteDomain(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "localhost", c.SiteDomain())
	c.options.SiteUrl = "https://foo.bar.com:2342/"
	assert.Equal(t, "foo.bar.com", c.SiteDomain())
	c.options.SiteUrl = ""
	assert.Equal(t, "localhost", c.SiteDomain())
}

func TestConfig_SitePreview(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.Equal(t, "https://i.photoprism.app/prism?cover=64&style=centered%20dark&caption=none&title=PhotoPrism", c.SitePreview())
	c.options.SitePreview = "http://preview.jpg"
	assert.Equal(t, "http://preview.jpg", c.SitePreview())
	c.options.SitePreview = "preview123.jpg"
	assert.Equal(t, "http://localhost:2342/preview123.jpg", c.SitePreview())
	c.options.SitePreview = "foo/preview123.jpg"
	assert.Equal(t, "http://localhost:2342/foo/preview123.jpg", c.SitePreview())
	c.options.SitePreview = "/foo/preview123.jpg"
	assert.Equal(t, "http://localhost:2342/foo/preview123.jpg", c.SitePreview())
}

func TestConfig_SiteTitle(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "PhotoPrism", c.SiteTitle())
	c.options.SiteTitle = "Cats"
	assert.Equal(t, "Cats", c.SiteTitle())
}

func TestConfig_Serial(t *testing.T) {
	c := NewConfig(CliTestContext())

	result := c.Serial()

	t.Logf("Serial: %s", result)

	assert.NotEmpty(t, result)
}

func TestConfig_SerialChecksum(t *testing.T) {
	c := NewConfig(CliTestContext())

	serial := "zr2g80wvjmm1zwzg"
	expected := "c7dcdb1c"

	c.serial = serial

	result := c.SerialChecksum()

	// t.Logf("Serial: %s", c.serial)
	// t.Logf("SerialChecksum: %s", result)

	assert.NotEmpty(t, result)
	assert.Equal(t, expected, result)
}

func TestConfig_Public(t *testing.T) {
	c := NewConfig(CliTestContext())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "public"
	assert.True(t, c.Public())
	c.options.Demo = true
	c.options.Public = false
	c.options.AuthMode = "public"
	assert.True(t, c.Public())
	c.options.Demo = true
	c.options.Public = true
	c.options.AuthMode = "public"
	assert.True(t, c.Public())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "other"
	assert.False(t, c.Public())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "password"
	assert.False(t, c.Public())
	c.options.Demo = false
	c.options.Public = true
	c.options.AuthMode = "password"
	assert.True(t, c.Public())
	c.options.Demo = true
	c.options.Public = false
	c.options.AuthMode = "password"
	assert.True(t, c.Public())
}

func TestConfig_Auth(t *testing.T) {
	c := NewConfig(CliTestContext())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "public"
	assert.False(t, c.Auth())
	c.options.Demo = true
	c.options.Public = false
	c.options.AuthMode = "public"
	assert.False(t, c.Auth())
	c.options.Demo = true
	c.options.Public = true
	c.options.AuthMode = "public"
	assert.False(t, c.Auth())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "other"
	assert.True(t, c.Auth())
	c.options.Demo = false
	c.options.Public = false
	c.options.AuthMode = "password"
	assert.True(t, c.Auth())
	c.options.Demo = false
	c.options.Public = true
	c.options.AuthMode = "password"
	assert.False(t, c.Auth())
	c.options.Demo = true
	c.options.Public = false
	c.options.AuthMode = "password"
	assert.False(t, c.Auth())
}

func TestConfigOptions(t *testing.T) {
	c := NewConfig(CliTestContext())
	r := c.Options()
	assert.False(t, r.DisableExifTool)
	assert.Equal(t, r.AutoImport, 7200)
	assert.Equal(t, r.AutoIndex, -1)

	c.options = nil
	r2 := c.Options()
	assert.Equal(t, r2.AutoImport, 0)
	assert.Equal(t, r2.AutoIndex, 0)
}
