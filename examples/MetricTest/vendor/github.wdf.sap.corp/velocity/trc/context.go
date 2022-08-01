// (C) 2016-2021 SAP SE or an SAP affiliate company. All rights reserved.
package trc

import (
	"bytes"
	"context"
)

type ctxKey int

const (
	itemKey ctxKey = iota
	userKey
)

func checkIfKeyExists(ctx context.Context, parentFrags []infoFrag, info Info) (bool, context.Context) {
	for i, cur := range parentFrags {
		index := bytes.IndexByte(cur, '=')
		if index != len(info.K) {
			continue
		}

		if string(cur[:index]) == info.K {
			var buf bytes.Buffer
			buf.Write(cur[:index+1]) // including '='
			serializeInfoValue(&buf, info, nil)

			subFrags := make(infoFrags, len(parentFrags))
			copy(subFrags, parentFrags)
			subFrags[i] = buf.Bytes()

			return true, context.WithValue(ctx, itemKey, subFrags)
		}
	}
	return false, nil
}

// UpdateTrcInfo returns a new context updated with the key passed in info.
// If the key doesn't exist, proceeds like AttachTrcInfo
func UpdateTrcInfo(ctx context.Context, info Info) context.Context {
	parentFrags := infoFragmentsFromContext(ctx)
	ok, updatedCtx := checkIfKeyExists(ctx, parentFrags, info)
	if ok {
		return updatedCtx
	}
	subFrags := parentFrags.subFromFrag(serializeInfos(info))
	return context.WithValue(ctx, itemKey, subFrags)
}

// AttachTrcInfo returns a new context witch the given ContextItems attached in a serialized form
func AttachTrcInfo(ctx context.Context, info ...Info) context.Context {
	parentFrags := infoFragmentsFromContext(ctx)
	subFrags := parentFrags.subFromFrag(serializeInfos(info...))
	return context.WithValue(ctx, itemKey, subFrags)
}

func infoFragmentsFromContext(ctx context.Context) infoFrags {
	frags, _ := ctx.Value(itemKey).(infoFrags)
	return frags
}
