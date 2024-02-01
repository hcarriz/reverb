package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("invalid hash")
	ErrIncompatibleVariant = errors.New("incompatible variant")
	ErrIncompatibleVersion = errors.New("incompatible version")
)

const Variant = "argon2id"

type config struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  int
	keyLength   uint32
}

type Option interface {
	apply(*config) error
}

type option func(*config) error

func (o option) apply(c *config) error {
	return o(c)
}

func MemoryInMB(m uint32) Option {
	return Memory(m * 1024)
}

func Memory(mem uint32) Option {
	return option(func(c *config) error {
		c.memory = mem
		return nil
	})
}

func Iterations(it uint32) Option {
	return option(func(c *config) error {
		c.iterations = it
		return nil
	})
}

func Parallelism(p uint8) Option {
	return option(func(c *config) error {
		c.parallelism = p
		return nil
	})
}

func SaltLength(length int) Option {
	return option(func(c *config) error {
		c.saltLength = length
		return nil
	})
}

func KeyLength(length uint32) Option {
	return option(func(c *config) error {
		return nil
	})
}

type Password struct {
	Variant     string
	Version     int
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	Salt        []byte
	Key         []byte
}

func (p Password) String() string {

	b64Salt := base64.RawStdEncoding.EncodeToString(p.Salt)
	b64Key := base64.RawStdEncoding.EncodeToString(p.Key)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", p.Variant, p.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Key)
}

// Create a password.
func Create(pwd string, opts ...Option) (*Password, error) {

	var (
		err error
		c   = config{
			memory:      64 * 1024,
			iterations:  1,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		}
	)

	for _, opt := range opts {
		err = errors.Join(err, opt.apply(&c))
	}

	if err != nil {
		return nil, err
	}

	p := &Password{
		Variant:     Variant,
		Version:     argon2.Version,
		Memory:      c.memory,
		Iterations:  c.iterations,
		Parallelism: c.parallelism,
	}

	p.Salt, err = generatedBytes(c.saltLength)
	if err != nil {
		return nil, err
	}

	p.Key = argon2.IDKey([]byte(pwd), p.Salt, p.Iterations, p.Memory, p.Parallelism, uint32(len(p.Salt)))

	return p, nil
}

// Check that the plaintext and hashed string match.
func Check(plaintext, hashed string) (bool, error) {

	pwd, err := Parse(hashed)
	if err != nil {
		return false, err
	}

	return Compare(plaintext, *pwd)

}

// Compare a plaintext string and the Password type. Return their results or any errors.
func Compare(plaintext string, pwd Password) (bool, error) {

	maybe := argon2.IDKey([]byte(plaintext), pwd.Salt, pwd.Iterations, pwd.Memory, pwd.Parallelism, uint32(len(pwd.Salt)))

	if subtle.ConstantTimeCompare(pwd.Key, maybe) == 1 {
		return true, nil
	}

	return false, nil
}

func generatedBytes(amount int) ([]byte, error) {
	b := make([]byte, amount)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Parse will take a string attempt to parse a Password from it.
func Parse(hash string) (*Password, error) {

	vals := strings.Split(hash, "$")
	var err error

	if len(vals) != 6 {
		return nil, ErrInvalidHash
	}

	if vals[1] != Variant {
		return nil, ErrIncompatibleVariant
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, err
	}

	if version != argon2.Version {
		return nil, ErrIncompatibleVersion
	}

	params := &Password{}
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism); err != nil {
		return nil, err
	}

	params.Salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, err
	}

	params.Key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, err
	}

	return params, nil
}
