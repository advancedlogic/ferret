package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func GetInDocument(ctx context.Context, doc drivers.HTMLDocument, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "url", "URL":
			return doc.GetURL(), nil
		case "body":
			return doc.QuerySelector(ctx, "body"), nil
		case "head":
			return doc.QuerySelector(ctx, "head"), nil
		default:
			return GetInNode(ctx, doc.DocumentElement(), path)
		}
	}

	return GetInNode(ctx, doc.DocumentElement(), path)
}

func GetInElement(ctx context.Context, el drivers.HTMLElement, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "innerText":
			return el.InnerText(ctx), nil
		case "innerHTML":
			return el.InnerHTML(ctx), nil
		case "value":
			return el.GetValue(ctx), nil
		case "attributes":
			return el.GetAttributes(ctx), nil
		default:
			return GetInNode(ctx, el, path)
		}
	}

	return GetInNode(ctx, el, path)
}

func GetInNode(ctx context.Context, node drivers.HTMLNode, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	nt := node.Type()
	segment := path[0]
	st := segment.Type()

	if st == types.Int {
		if nt == drivers.HTMLElementType || nt == drivers.HTMLDocumentType {
			re := node.(drivers.HTMLNode)

			return re.GetChildNode(ctx, segment.(values.Int)), nil
		}

		return values.GetIn(ctx, node, path[0:])
	}

	if st == types.String {
		segment := segment.(values.String)

		switch segment {
		case "nodeType":
			return node.NodeType(), nil
		case "nodeName":
			return node.NodeName(), nil
		case "children":
			return node.GetChildNodes(ctx), nil
		case "length":
			return node.Length(), nil
		default:
			return values.None, nil
		}
	}

	return values.None, core.TypeError(st, types.Int, types.String)
}
