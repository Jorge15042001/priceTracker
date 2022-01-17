import requests
res = requests.post('http://localhost:5000/', json={"url":"https://www.amazon.com/LG-32GN650-B-Ultragear-Reduction-FreeSync/dp/B08LLF9NS1/ref=psdc_1292115011_t4_B07YGZRQ98"})
if res.ok:
    print(res.json())
