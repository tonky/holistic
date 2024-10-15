package {{ domain }}

{% for i in imports %}
import "{{ i }}"
{% end %}

type {{ model }} struct {
    {% for field in fields %}
        {{ field.Name }} {{ field.T.StructType(domain) }}
    {% end %}

}