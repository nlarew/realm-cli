.. _{{.CommandRef}}:

{{range toChars .CommandName}}={{end}}
{{.CommandName}}
{{range toChars .CommandName}}={{end}}

.. program:: {{.CommandName}}

{{.Description}}

.. code-block:: bash
   
   {{.Usage}}

.. _{{.CommandRef}}-options:

Options
-------

.. list-table::
   :header-rows: 1
   :width: 55 5 40
   
   * - Option
     - Shorthand
     - Description
   {{range $index, $option := .Options}}
   * - .. option:: --{{$option.Name}}
     - {{if $option.Shorthand}}``-{{$option.Shorthand}}``{{end}}
     - {{if $option.DefaultValue}}Default: {{$option.DefaultValue}}
       
       {{$option.Usage}}{{else}}{{$option.Usage}}{{end}}
{{end}}
.. _{{.CommandRef}}-inherited-options:

Inherited Options
-----------------

.. list-table::
   :header-rows: 1
   :width: 55 5 40
   
   * - Option
     - Shorthand
     - Description
   {{range $index, $option := .InheritedOptions}}
   * - .. option:: --{{$option.Name}}
     - {{if $option.Shorthand}}``-{{$option.Shorthand}}``{{end}}
     - {{if $option.DefaultValue}}Default: {{$option.DefaultValue}}
       
       {{$option.Usage}}{{else}}{{$option.Usage}}{{end}}
{{end}}
{{if .SeeAlso}}
See Also
--------
{{range .SeeAlso}}
* {{.}}
{{end}}
{{end}}
