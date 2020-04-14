# dyn-connector-go
Go M2M interface to Dyn/PigrecoOS

## Features
* Stateless and stateful modes
* Transparent Auth and tokens expiration management
 
## Installation and usage
```
import "github.com/modulo-srl/dyn-connector-go
```

See `connector_test.go` for examples.

**Master Token** and **Auth UID** are authorization constants that could be obtained from [PigrecoOS](http://www.pigrecoos.it) services.

**Session token**, used with the only get/set callback required by the Connector, must to be saved on a persistent storage, for example a file or a DB.

For **Operation** and related parameters please refer to [PigrecoOS](http://www.pigrecoos.it) API section.

---

*Copyright 2020 [Modulo srl](http://www.modulo.srl) - Licensed under the Apache license*
