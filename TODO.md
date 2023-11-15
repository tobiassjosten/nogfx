- Byt namn på pkg/ till internal/.

- Byt till det officiella Simpex-libbet.

- Bryt ut telnet till ett externt bibliotek.
	- Gör om Commands från en kanal till en custom callback-funktion.
	- Använd callback för negotiation, med en färdig funktion (som går att chaina).

- Testa fx för dependecy injection i main().

## gmcp

- Lägg till tester för data av fel typer. `[]` för objekt, `{}` för listor, `""` för nummer, etc.

- Lägg till tester för avsaknad av data, t.ex. `Char.Items.Contents` endast.


## scripting

js:
https://github.com/robertkrimen/otto
https://github.com/dop251/goja

lua:
https://github.com/milochristiansen/lua
https://github.com/yuin/gopher-lua
https://github.com/Shopify/go-lua

other:
https://github.com/d5/tengo
https://github.com/ozanh/ugo

	script := tengo.NewScript([]byte(`asdf("lololol")`))
	err := script.Add(
		"asdf",
		func(args ...tengo.Object) (tengo.Object, error) {
			for _, arg := range args {
				world.ui.Print([]byte(arg.String()))
			}
			return tengo.FromInterface("zxcv")
		},
	)
	if err != nil {
		log.Println("failed adding", err)
	}
	_, err = script.Run()
	if err != nil {
		log.Println("failed running", err)
	}

- Döp om t.ex. Inoutput.RemoveInput till Inoutput.InputRemove (för att bättre
  signalera relationen till Inoutput.Input.Remove).

- Bör TUI också använda pkg.Inoutput?

- Kan vi ta bort kanalerna från TUI och Client? Och lägga det i Engine istället?

- Client skulle kunna ha en "DefaultCommandee" som kastar bort commands när de
  kommer in. Om man skickar in sin egen custom-Commandee så får den istället ta
  emot commands. Då slipper vi exponera en kanal från Client.

- Ersätt alla nogfx.org med nogfx.com.

- IOKind känns som ett hack. Det måste finnas ett mer elegant sätt att lösa
  det på.
