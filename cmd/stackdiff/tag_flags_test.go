package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTagFlags_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseTagFlags(fs, []string{})
	require.NoError(t, err)
	assert.Empty(t, cfg.Rules)
}

func TestParseTagFlags_SingleTypeRule(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseTagFlags(fs, []string{"--tag", "type=aws_instance:compute"})
	require.NoError(t, err)
	require.Len(t, cfg.Rules, 1)
	assert.Equal(t, "aws_instance", cfg.Rules[0].ResourceType)
	assert.Equal(t, "compute", cfg.Rules[0].Tag)
}

func TestParseTagFlags_SinglePrefixRule(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseTagFlags(fs, []string{"--tag", "prefix=module.network:networking"})
	require.NoError(t, err)
	require.Len(t, cfg.Rules, 1)
	assert.Equal(t, "module.network", cfg.Rules[0].AddressPrefix)
	assert.Equal(t, "networking", cfg.Rules[0].Tag)
}

func TestParseTagFlags_MultipleRules(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseTagFlags(fs, []string{"--tag", "type=aws_s3_bucket:storage,prefix=module.app:frontend"})
	require.NoError(t, err)
	assert.Len(t, cfg.Rules, 2)
}

func TestParseTagFlags_InvalidRule_Skipped(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg, err := ParseTagFlags(fs, []string{"--tag", "badformat"})
	require.NoError(t, err)
	assert.Empty(t, cfg.Rules)
}

func TestParseTagRule_UnknownKey(t *testing.T) {
	_, ok := parseTagRule("unknown=foo:bar")
	assert.False(t, ok)
}
