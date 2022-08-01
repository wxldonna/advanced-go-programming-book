package trc

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey int

const (
	itemKey ctxKey = iota
	userKey
)

func checkIfKeyExists(ctx context.Context, parentFrags []infoFrag, info Info) (bool, context.Context) {
	for i, cur := range parentFrags {
		content := string(cur[:])
		key := strings.Split(content, "=")[0]
		if key == info.K {
			newValue := fmt.Sprintf("%s=\"%s\"", key, info.V)
			infoInBytes := []byte(newValue)
			subFrags := make(infoFrags, len(parentFrags))
			copy(subFrags, parentFrags)
			subFrags[i] = infoInBytes
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
