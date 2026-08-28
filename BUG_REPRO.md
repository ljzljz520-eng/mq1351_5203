# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	qin-culture-site/cmd/qinweb	0.013s
ok  	qin-culture-site/internal/catalog	0.013s
ok  	qin-culture-site/internal/config	0.014s
ok  	qin-culture-site/internal/domain	0.013s
--- FAIL: TestMusicTitleFollowsSelection (0.00s)
    music_test.go:26: title should follow current selection, got "流水"
FAIL
FAIL	qin-culture-site/internal/music	0.012s
ok  	qin-culture-site/internal/service	0.009s
ok  	qin-culture-site/internal/store	0.023s
ok  	qin-culture-site/internal/web	0.009s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/qinweb): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/qinweb): exit `0`
