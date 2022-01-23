import requests
res = requests.post('http://localhost:5000/', json={"url":"https://articulo.mercadolibre.com.ec/MEC-506919172-apple-ipad-pro-11-128gb-wifi-chip-m1-_JM?searchVariation=174101054149#searchVariation=174101054149&position=4&search_layout=grid&type=item&tracking_id=105f8369-73e0-4076-ba01-c8d78b96e1c8"})
if res.ok:
    print(res.json())
