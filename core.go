package debug

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Flag int

const (
	TimeLocal Flag = 1 << iota
	TimeUTC   Flag = 1 << iota
	LowerCase Flag = 1 << iota
	UpperCase Flag = 1 << iota
	Mills     Flag = 1 << iota
)

var (
	colors           = make([]*color.Color, 0)
	allows           = make([]func(string) bool, 0)
	noop             = &noopDebug{}
	locker           = &sync.RWMutex{}
	warehouse        = make(map[uint]IDebug)
	defaultFlag Flag = 0
)

func init() {
	envDebug := os.Getenv("DEBUG")
	validator := regexp.MustCompile("([a-zA-Z0-9_\\-]+|\\*)(:[a-zA-Z0-9_\\-]+)*(:\\*)?")
	rep := strings.NewReplacer("*", ".+")
	for _, it := range strings.Split(envDebug, ",") {
		str := strings.TrimSpace(it)
		if !validator.MatchString(str) {
			continue
		}
		va := regexp.MustCompile(rep.Replace(str))
		allows = append(allows, func(s string) bool {
			return va.MatchString(s)
		})
	}
	for _, it := range []color.Attribute{
		color.FgRed,
		color.FgHiRed,
		color.FgGreen,
		color.FgHiGreen,
		color.FgYellow,
		color.FgHiYellow,
		color.FgBlue,
		color.FgHiBlue,
		color.FgMagenta,
		color.FgHiMagenta,
		color.FgCyan,
		color.FgHiCyan,
		color.FgWhite,
		color.FgHiWhite,
	} {
		co := color.New(it)
		co.EnableColor()
		colors = append(colors, co)
	}
}

type IDebug interface {
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

type myDebug struct {
	nsp    string
	c      *color.Color
	f      Flag
	t      time.Time
	locker *sync.Mutex
}

func (p *myDebug) getPrefix() string {
	if p.f&TimeLocal == TimeLocal {
		return p.c.Sprintf("%s [%s]:", time.Now().Format("2006-01-02T15:04:05.000Z07:00"), p.nsp)
	} else if p.f&TimeUTC == TimeUTC {
		return p.c.Sprintf("%s [%s]:", time.Now().In(time.UTC).Format("2006-01-02T15:04:05.000Z"), p.nsp)
	}
	return p.c.Sprintf("[%s]:", p.nsp)
}

func (p *myDebug) getSuffix() *string {
	if p.f&Mills != Mills {
		return nil
	}
	var old time.Time
	old, p.t = p.t, time.Now()
	nano := p.t.UnixNano() - old.UnixNano()
	s := p.c.Sprintf("+%dms", nano/1000000)
	return &s
}

func (p *myDebug) Print(a ...interface{}) {
	p.locker.Lock()
	defer p.locker.Unlock()
	fmt.Fprint(os.Stdout, p.getPrefix(), " ")
	fmt.Fprint(os.Stdout, a...)
	var suffix = p.getSuffix()
	if suffix != nil {
		fmt.Fprint(os.Stdout, " ", *suffix)
	}
}

func (p *myDebug) Println(a ...interface{}) {
	p.locker.Lock()
	defer p.locker.Unlock()
	fmt.Fprint(os.Stdout, p.getPrefix(), " ")
	var suffix = p.getSuffix()
	if suffix != nil {
		fmt.Fprintln(os.Stdout, append(a, *suffix)...)
	} else {
		fmt.Fprintln(os.Stdout, a...)
	}
}

func (p *myDebug) Printf(format string, a ...interface{}) {
	p.locker.Lock()
	defer p.locker.Unlock()
	var suffix = p.getSuffix()
	if suffix == nil {
		fmt.Fprintf(os.Stdout, p.getPrefix()+" "+format, a...)
		return
	}
	if strings.HasSuffix(format, "\n") {
		fmt.Fprintf(os.Stdout, p.getPrefix()+" "+strings.TrimRight(format, "\n"), a...)
		fmt.Fprint(os.Stdout, " ", *suffix, "\n")
	} else {
		fmt.Fprintf(os.Stdout, p.getPrefix()+" "+format, a...)
		fmt.Fprint(os.Stdout, " ", *suffix)
	}
}

func SetFlags(first Flag, others ...Flag) {
	defaultFlag |= first
	for _, it := range others {
		defaultFlag |= it
	}
}

func Debug(namespace string, flags ...Flag) IDebug {
	// check is debug allow
	var ok bool
	for _, it := range allows {
		ok = it(namespace)
		if ok {
			break
		}
	}
	// debug is disabled,return noop debug.
	if !ok {
		return noop
	}
	var f = defaultFlag
	for _, it := range flags {
		f |= it
	}
	var h = hashcode(fmt.Sprintf("%s@%d", namespace, f))
	locker.RLock()
	found, ok := warehouse[h]
	locker.RUnlock()
	if ok {
		return found
	}
	var h2 = hashcode(namespace) % uint(len(colors))
	co := colors[int(h2)]
	locker.Lock()
	var ret IDebug
	var nsp = namespace
	if f&LowerCase == LowerCase {
		nsp = strings.ToLower(namespace)
	} else if f&UpperCase == UpperCase {
		nsp = strings.ToUpper(namespace)
	}
	ret = &myDebug{nsp, co, f, time.Now(), &sync.Mutex{}}
	warehouse[h] = ret
	locker.Unlock()
	return ret
}

func hashcode(s string) uint {
	var h uint
	for _, b := range []byte(s) {
		h = h<<5 - h + uint(b)
	}
	return h
}

type noopDebug struct {
}

func (p *noopDebug) Print(a ...interface{}) {
}

func (p *noopDebug) Println(a ...interface{}) {
}

func (p *noopDebug) Printf(format string, a ...interface{}) {
}
