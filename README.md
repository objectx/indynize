# indynize

Replace JARs in `$GROOVY_HOME/lib` with *InvokeDynamic* enabled versions.

Performed operations are basically described [here](http://groovy-lang.org/indy.html).
But the example using `sh` and other UN?Xish utilities. It's not feasible for Windows environment.
This is the reason I wrote this program (1st in **Python**, now using **Go**).

Original JARs are kept in `$GROOVY_HOME/lib.orig`, so you can rollback to original state, just doing:
```bash
$ (cd $GROOVY_HOME ; rm -rf lib ; mv lib.orig lib)
```

### LICENSE: [MIT][]

[MIT]: http://opensource.org/licenses/MIT
