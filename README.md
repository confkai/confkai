# confkai

Confkai is a configuration as code library for Go. With this library
you compose your configurations with functions. This no-frills
library includes the basic composition functions to get you started. But,
you can write your own functions or import a separate module to support values
from other providers like GCP Secrets or AWS ParamStore! 

Benefits of using Confkai for configuration:   
1. Changes are tracked by git.    
1. Changes can be audited in every PR.  
1. Roll back Safe deployments.  
1. Lazy Loading by default.  
1. A single source of truth.  
1. Eager loading, Caching, and much more included with the base library.  
1. No dependencies in the base library.  


