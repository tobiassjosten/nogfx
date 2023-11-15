package procs

/* @todo Move to its own module.
var (
	bal = time.Time{}
	eq  = time.Time{}
)

// Requires: prompt *hh*1, *mm*2, *ee*3, *ww*4 *Rr *rk *b*c*d *s
ps1 := []byte("\x1b[32m4")
ps2 := []byte("\x1b[37m\x1b[32m4")
if (len(output) > len(ps1) && bytes.Equal(output[:len(ps1)], ps1)) ||
	(len(output) > len(ps2) && bytes.Equal(output[:len(ps2)], ps2)) {
	loutput := len(output)

	xstamp := "15:04:05.000"
	lstamp := len(xstamp) - 1
	sstamp := string(output[loutput-lstamp-1:loutput-1]) + "0"
	tstamp, _ := time.Parse(xstamp, sstamp)

	prlbal := output[loutput-lstamp-4] == 'R'
	pllbal := output[loutput-lstamp-5] == 'L'
	pbal := output[loutput-lstamp-6] == 'x'

	eqoffset := 0
	if !pbal {
		eqoffset = 1
	}
	peq := output[loutput-lstamp-7+eqoffset] == 'e'
	pbal = pbal && prlbal && pllbal

	if pbal && bal != (time.Time{}) {
		diff := fmt.Sprintf("\x1b[30;1m %.2fx\x1b[37m", tstamp.Sub(bal).Seconds())
		output = append(output, []byte(diff)...)
		bal = time.Time{}
	} else if !pbal && bal == (time.Time{}) {
		bal = tstamp
	}

	if peq && eq != (time.Time{}) {
		diff := fmt.Sprintf("\x1b[30;1m %.2fe\x1b[37m", tstamp.Sub(eq).Seconds())
		output = append(output, []byte(diff)...)
		eq = time.Time{}
	} else if !peq && eq == (time.Time{}) {
		eq = tstamp
	}
}

KOM IHÅG att ta höjd för blackout. Vi skulle kunna tracka bal/eq i prompten och jämföra med föregående för att se ifall vi förlorat bal/eq sedan sist (eller om det var okänt, som med blackout, så skiter vi i att tajma).
*/
