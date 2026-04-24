//ff:func feature=pkg-auth type=util control=sequence topic=auth-refresh
//ff:what marshalJSON — json.Marshal 래퍼 (external import 간결화)
package auth

import "encoding/json"

func marshalJSON(v any) ([]byte, error) {
	if raw, ok := v.(json.RawMessage); ok {
		return raw, nil
	}
	return json.Marshal(v)
}
