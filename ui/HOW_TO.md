###First check it exists a version of your library compatible with the version of Angular defined in package.json.
How to install ng-bootstrap

ng add @ng-bootstrap/ng-bootstrap

How to install ngx-bootstrap

ng add ngx-bootstrap

How to install Angular Material

Add Angular Material using ng add command:

ng add @angular/material

You will get the following questions:

? Choose a prebuilt theme name, or "custom" for a custom theme: *Choose any theme you like here*
? Set up global Angular Material typography styles? *Yes* 
? Set up browser animations for Angular Material? *Yes*

Angular Material will start installing, but you will get the following error after installation:

Your project is not using the default builders for "build". The Angular Material schematics cannot add a theme to the workspace configuration if the builder has been changed.

No need to Panic! Just add your desired theme in style.scss:

import '@angular/material/prebuilt-themes/indigo-pink.css'

Angular Material Library is now installed in your project.
