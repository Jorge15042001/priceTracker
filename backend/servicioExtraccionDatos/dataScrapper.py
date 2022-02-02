from selenium.webdriver import Chrome
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
from abc import ABC, abstractmethod
from tldextract import extract
import json
from flask import Flask,request
from datetime import date
app = Flask(__name__)


chrome_options = Options()
chrome_options.add_argument("--headless")


months = [ 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov','Dec']
months_sp = [ 'ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov','Dic']
def getMonthNumber(monthName:str) ->int:
    abv = monthName[:3]
    print("abv",abv)
    for i in range (len(months)):
        if months[i].lower()==abv.lower():
            return i + 1
    for i in range (len(months_sp)):
        if months_sp[i].lower()==abv.lower():
            return i + 1
    raise Exception("not valid month name")

class productInformation:
    link :str
    imgSrc : str
    opinion : float
    arrivalDate : str 
    price : float

    def __init__(self,link:str,img:str,opinion:float,arrivalDate:str,price:float):
        self.link=link
        self.imgSrc = img
        self.opinion =opinion
        self.arrivalDate = arrivalDate
        self.price = price

class WebStoreScraper(ABC):
    @abstractmethod
    def execute(self,driver) -> productInformation:
        pass

class AmazonScraper(WebStoreScraper):
    def execute(self,driver) -> productInformation:
        price = driver.find_element(By.CLASS_NAME,"a-price-whole").text
        image_src = driver.find_element(By.ID,"landingImage").get_attribute("src")
        puntuacion = driver.find_element(By.CSS_SELECTOR, "#acrPopover .a-declarative .a-popover-trigger .a-icon .a-icon-alt").get_attribute("innerText").split()[0]
        date_fields = driver.find_element(By.ID,"mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE").text.split(",")[1].split(".")[0].strip(" ").split(" ")
        print(date_fields)
        month = getMonthNumber(date_fields[0]) 
        day = int(date_fields[1])
        arrival_date = date(2022,month,day)

        return productInformation("",image_src,float(puntuacion),str(arrival_date),float(price))
        # return {"precio":float(price),"image_src":image_src,"puntuacion":float(puntuacion),"fecha_entrega":str(arrival_date)}

class EbayScraper(WebStoreScraper):
    def execute(self,driver) -> productInformation:
        price = driver.find_element(By.ID,"prcIsum").text.split()[1][1:]
        image_src = driver.find_element(By.ID,"icImg").get_attribute("src")
        puntuacion = driver.find_element(By.CSS_SELECTOR, "#review-ratings-cntr .reviews-star-rating").get_attribute("title").split(",")[0].split()[0]
        arrival_date_fields = driver.find_element(By.CSS_SELECTOR,"#delSummary .sh-del-frst .vi-acc-del-range ").text.split(".")[1].strip(" ").split(" ")
        arrival_day = int(arrival_date_fields[0])
        arrival_month = getMonthNumber(arrival_date_fields[1])

        arrival_date = date(2022,arrival_month,arrival_day)
        return productInformation("",image_src,float(puntuacion),str(arrival_date),float(price))
        # return {"precio":float(price),"image_src":image_src,"puntuacion":float(puntuacion),"fecha_entrega":str(arrival_date)}

class MercadoLibreScraper(WebStoreScraper):
    def execute(self,driver) -> productInformation:
        price = "".join(driver.find_element(By.CLASS_NAME,"price-tag-amount").text.split("\n")[1:]).replace(",",".")
        price_segments = price.split(".")
        if len(price_segments[-1])==3:
            price = "".join(price_segments)
        print("\n\n",price,"\n\n")
        image_src = driver.find_element(By.CSS_SELECTOR,".ui-pdp-gallery__figure img").get_attribute("src")
        puntuacion = driver.find_element(By.CLASS_NAME, "ui-pdp-seller__sales-description").text.strip("%")
        arrival_date = date.today()
        return productInformation("",image_src,float(puntuacion)/20,str(arrival_date),float(price))
        # return {"precio":float(price),"image_src":image_src,"puntuacion":float(puntuacion)/20,"fecha_entrega":str(fecha_entrega)}


def getScrapingStrategy(url:str)  ->WebStoreScraper:
    _, td, tsu = extract(url) # prints abc, hostname, com
    domain = td + '.' + tsu.split(".")[0] # will prints as hostname.com    
    if domain == "amazon.com":return AmazonScraper()
    if domain == "ebay.com":return EbayScraper()
    if domain == "mercadolibre.com":return MercadoLibreScraper()
    raise Exception("No avaliable scraper")

class WebScraper_onlineStore:
    scrapingStrategy: WebStoreScraper  ## the strategy interface
    def __init__(self,url:str):
        self.driver =Chrome("/usr/bin/chromedriver",options=chrome_options)
        self.driver.get(url)

    def setScrapingMethod(self, strategy: WebStoreScraper ) -> None:
        self.scrapingStrategy = strategy

    def run(self) -> productInformation:
        return  self.scrapingStrategy.execute(self.driver)


@app.route('/',methods=["POST"])
def parse():
    url = str(request.form.get("url"))
    scraper = WebScraper_onlineStore(url)

    scrapingStrategy = getScrapingStrategy(url)
    scraper.setScrapingMethod(scrapingStrategy)

    info = scraper.run()
    info.link = url
    return json.dumps(info.__dict__)


app.run(debug=True)

driver.quit() # closing the browser

