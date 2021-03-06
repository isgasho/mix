package xfmt

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "regexp"
    "testing"
)

type level1 struct {
    Name     string
    Level2   *level2
    Level2_1 *level2
    Level4   *level4
    Level4_1 level4
}

type level2 struct {
    Name   string
    Level3 *level3
}

type level3 struct {
    Name string
}

type level4 struct {
    Name   string
    Level5 *level5
}

type level5 struct {
    Name string
}

type level6 struct {
    Name  string
    Map   map[string]*level5
    Array []*level5
}

func filterAddr(a string) string {
    reg, _ := regexp.Compile("0x[0-9a-z]*")
    return string(reg.ReplaceAll([]byte(a), []byte("")))
}

func TestSprintf(t *testing.T) {
    a := assert.New(t)

    l5 := level5{Name: "Level5"}
    l4 := level4{Name: "Level4", Level5: &l5}
    l3 := level3{Name: "Level3"}
    l2 := level2{Name: "Level2", Level3: &l3}
    l1 := level1{Name: "level1", Level2: &l2, Level2_1: &l2, Level4: &l4, Level4_1: l4}

    a.Equal(filterAddr(Sprintf(1, "%+v", l1)), filterAddr("{Name:level1 Level2:0xc00000c100 Level2_1:0xc00000c100 Level4:0xc00000c0e0 Level4_1:{Name:Level4 Level5:0xc0000404d0}}"))
    a.Equal(filterAddr(Sprintf(1, "%+v", &l1)), filterAddr("&{Name:level1 Level2:0xc00000c100 Level2_1:0xc00000c100 Level4:0xc00000c0e0 Level4_1:{Name:Level4 Level5:0xc0000404d0}}"))
    a.Equal(filterAddr(Sprintf(2, "%+v", l1)), filterAddr("{Name:level1 Level2:0xc00000c100:&{Name:Level2 Level3:0xc0000404e0} Level2_1:0xc00000c100 Level4:0xc00000c0e0:&{Name:Level4 Level5:0xc0000404d0} Level4_1:{Name:Level4 Level5:0xc0000404d0}}"))
    a.Equal(filterAddr(Sprintf(3, "%+v", l1)), filterAddr("{Name:level1 Level2:0xc00000c100:&{Name:Level2 Level3:0xc0000404e0:&{Name:Level3}} Level2_1:0xc00000c100 Level4:0xc00000c0e0:&{Name:Level4 Level5:0xc0000404d0:&{Name:Level5}} Level4_1:{Name:Level4 Level5:0xc0000404d0}}"))
    a.Equal(filterAddr(Sprintf(100, "%+v", l1)), filterAddr("{Name:level1 Level2:0xc00000c100:&{Name:Level2 Level3:0xc0000404e0:&{Name:Level3}} Level2_1:0xc00000c100 Level4:0xc00000c0e0:&{Name:Level4 Level5:0xc0000404d0:&{Name:Level5}} Level4_1:{Name:Level4 Level5:0xc0000404d0}}"))
}

func TestSprintfMultiple(t *testing.T) {
    a := assert.New(t)

    l5 := level5{Name: "Level5"}
    l4 := level4{Name: "Level4", Level5: &l5}
    l3 := level3{Name: "Level3"}

    a.Equal(filterAddr(Sprintf(100, "%v And %+v And %#v", l3, l4, &l5)), filterAddr(`{Level3} And {Name:Level4 Level5:0xc0000404d0:&{Name:Level5}} And &xfmt.level5{Name:"Level5"}`))
    a.Equal(filterAddr(Sprintf(100, "%v And %+v And %#v And %d And %s", l3, l4, &l5, 123, "456")), filterAddr(`{Level3} And {Name:Level4 Level5:0xc00007efd0:&{Name:Level5}} And &xfmt.level5{Name:"Level5"} And 123 And 456`))
    a.Equal(filterAddr(Sprintf(100, "%d And %s And %v And %+v And %#v", 123, "456", l3, l4, &l5)), filterAddr(`123 And 456 And {Level3} And {Name:Level4 Level5:0xc00007efd0:&{Name:Level5}} And &xfmt.level5{Name:"Level5"}`))
    a.Equal(filterAddr(Sprintf(100, "%d And %v And %+v And %s And %#v", 123, l3, l4, "456", &l5)), filterAddr(`123 And {Level3} And {Name:Level4 Level5:0xc00007efd0:&{Name:Level5}} And 456 And &xfmt.level5{Name:"Level5"}`))
}

func TestMoreFunc(t *testing.T) {
    l5 := level5{Name: "Level5"}
    l4 := level4{Name: "Level4", Level5: &l5}

    fmt.Println(l4, &l5)
    Print(2, l4, &l5)
    println("")
    Println(2, l4, &l5)
    Printf(2, "%v %v\n", l4, &l5)
    println(Sprint(2, l4, &l5))
    print(Sprintln(2, l4, &l5))

    //{Level4 0xc000048fe0} &{Level5}
    //{Level4 0xc000048fe0:&{Level5}} &{Level5}
    //{Level4 0xc000048fe0:&{Level5}} &{Level5}
    //{Level4 0xc000048fe0:&{Level5}} &{Level5}
    //{Level4 0xc000048fe0:&{Level5}} &{Level5}
    //{Level4 0xc000048fe0:&{Level5}} &{Level5}
}

func TestMap(t *testing.T) {
    a := assert.New(t)

    m := map[string]*level5{}
    m["foo"] = &level5{}
    m["bar"] = &level5{}
    a.Equal(filterAddr(Sprintf(2, "%v", m)), filterAddr("map[bar:0xc000080ff0:&{} foo:0xc000080fe0:&{}]"))
    a.Equal(filterAddr(Sprintf(2, "%v", &m)), filterAddr("&map[bar:0xc000080ff0:&{} foo:0xc000080fe0:&{}]"))

    m2 := map[string]level5{}
    m2["foo"] = level5{}
    m2["bar"] = level5{}
    a.Equal(filterAddr(Sprintf(2, "%v", m2)), filterAddr("map[bar:{} foo:{}]"))
    a.Equal(filterAddr(Sprintf(2, "%v", &m2)), filterAddr("&map[bar:{} foo:{}]"))
}

func TestArray(t *testing.T) {
    a := assert.New(t)

    a1 := []*level5{}
    a1 = append(a1, &level5{})
    a1 = append(a1, &level5{})
    a.Equal(filterAddr(Sprintf(2, "%v", a1)), filterAddr("[0xc000049000:&{} 0xc000049010:&{}]"))
    a.Equal(filterAddr(Sprintf(2, "%v", &a1)), filterAddr("&[0xc000049000:&{} 0xc000049010:&{}]"))

    a2 := []level5{}
    a2 = append(a2, level5{})
    a2 = append(a2, level5{})
    a.Equal(filterAddr(Sprintf(2, "%v", a2)), filterAddr("[{} {}]"))
    a.Equal(filterAddr(Sprintf(2, "%v", &a2)), filterAddr("&[{} {}]"))
}

func TestStructMapArray(t *testing.T) {
    a := assert.New(t)

    m := map[string]*level5{}
    m["foo"] = &level5{Name: "Level5"}
    m["bar"] = &level5{Name: "Level5"}

    a1 := []*level5{}
    a1 = append(a1, &level5{Name: "Level5"})
    a1 = append(a1, &level5{Name: "Level5"})

    x := level6{
        Name:  "level6",
        Map:   m,
        Array: a1,
    }
    a.Equal(filterAddr(Sprintf(2, "%v", x)), filterAddr("{level6 map[bar:0xc000049020:&{Level5} foo:0xc000049010:&{Level5}] [0xc000049030:&{Level5} 0xc000049040:&{Level5}]}"))
    a.Equal(filterAddr(Sprintf(2, "%v", &x)), filterAddr("&{level6 map[bar:0xc000049020:&{Level5} foo:0xc000049010:&{Level5}] [0xc000049030:&{Level5} 0xc000049040:&{Level5}]}"))
}
