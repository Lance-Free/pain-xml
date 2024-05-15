# PainXML

The `painxml` package facilitates the generation and parsing of XML documents conforming to the ISO 20022 standard for
payment initiation, specifically focusing on Direct Debit documents. It offers structures and methods to convert payment
orders represented as Go structs into XML documents and vice versa.

## Overview

The package provides functionality to convert payment orders represented as Go structs into XML documents conforming to
the ISO 20022 standard. Additionally, it offers the capability to parse XML documents back into Go structs, allowing
seamless integration with existing systems that handle payment orders.

## Features

- **Conversion to XML**: Easily convert payment orders represented as Go structs into ISO 20022-compliant XML documents.
- **Parsing from XML**: Parse XML documents conforming to the ISO 20022 standard back into Go structs representing
  payment orders.

## Usage

### Generating Direct Debit Documents

To generate a Direct Debit XML document from a payment order represented as a Go struct, use the `ToDocument` function
provided by the `painxml` package:

```go
order := order.Order{
ExecutionDate: time.Now(),
// Populate order fields...
}

document, err := painxml.ToDocument(order)
if err != nil {
// Handle error
}
// Use generated XML document...
```

### Parsing Direct Debit Documents

To parse a Direct Debit XML document back into a payment order represented as a Go struct, use the ToOrder function
provided by the painxml package:

```go
xmlData := []byte("<xml>...</xml>") // Replace with actual XML data
var doc document.DirectDebit
err := xml.Unmarshal(xmlData, &document)
if err != nil {
// Handle error
}

order, err := painxml.ToOrder(doc)
if err != nil {
// Handle error
}
// Use parsed payment order...
```

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.