# plugin.go
Plugins in GOlang. It works, but it's rough around the edges as it's just a prototype

## Screencast
[![Screencast](http://img.youtube.com/vi/ZxtV3rRJssI/2.jpg)] (https://www.youtube.com/watch?v=ZxtV3rRJssI)

## Why?
One of the big things I love about interpreted languages is the extensibility of the appliation without needing to recompile it everytime. Although I'm not a huge fan of WordPress, one of the ideas/concepts that I really love that they implemented is hooks/filters. I wanted this ability in golang. 

## How?
It's really simple especially if you're familiar with hooks/filters in wordpress. Here is a quick bit of sample code: 

```
plugin.Init("/tmp", os.Args)
my_password := "asd123"
password := plugin.Filter("filter_password", my_password)
fmt.Println("Before: ", "asd123")
fmt.Println("Hashed: ", password)
```

Each plugin(script/application) you build will simply need two things. How to register itself, and reading in stdin. Plugin.go will pass your script the flag ```--register-plugin``` and this function _should_ return a json representation with an object inside. Here is a sample of it's representation: ```[{"action": "filter", "priority": 1, "trigger": "filter_password"}]``` 

```action``` should either be ```hook``` or ```filter```.
```priority``` should be an `int`.
```trigger``` should be a string of your choosing(you'll call this in hook/filter call in go)

Your application should then be capable of reading stdin, and spitting out it's modifications. 

## TL;DR
Write a script(any language) that can read stdin, and also can spit out json. plugin.go will call it, you mutate/do something and the output(stdout) is used.
