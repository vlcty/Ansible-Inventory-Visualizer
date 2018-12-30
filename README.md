# Ansible inventory visualizer

Render a table based on your ansible inventory. For example:

```
       SERVER       | DBSERVERS | WEBSERVERS |
+-------------------+-----------+------------+--+
  bar.example.com   |           | X          |
  foo.example.com   |           | X          |
  one.example.com   | X         |            |
  three.example.com | X         |            |
  two.example.com   | X         |            |
```

## Installation and usage

This project has two dependencies: [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) and [go-ini/ini](https://github.com/go-ini/ini). Install the dependencies first:

```
:-$ go get github.com/olekukonko/tablewriter
:-$ go get github.com/go-ini/ini
```

Then you can build it:

```
:-$ go build InventoryVisualizer.go
:-$ ./InventoryVisualizer --inventory test_inventories/simple_valid_inventory
```

That's it.

## Filter for specific hosts

You can filter for specific hosts passing a comma seperated list. For example:

```
:-$ ./InventoryVisualizer --inventory test_inventories/simple_valid_inventory --filter one.example.com,two.example.com
```

The output:

```
      SERVER      | DBSERVERS | WEBSERVERS |
+-----------------+-----------+------------+--+
  one.example.com | X         |            |
  two.example.com | X         |            |
```

## But ... why?

Good question! While developing this application I thought it was a good idea. Then I ran it against some inventories I found online (most of them from the official documentation) and they came up nice. Then I tested it agains the one from my work and it basically blew up. Too much columns and rows. You've been warned :-)

## Use a pager

Based on your inventory this application can produce a big output. Pipe the output into `less -S` to have a horizontal and vertical pager.
