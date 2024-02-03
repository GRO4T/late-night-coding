import logging

from selenium import webdriver
from selenium.webdriver.common.by import By

# ---------------------------------------------------------------------------- #
#                                    Logging                                   #
# ---------------------------------------------------------------------------- #
format = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
logging.basicConfig(level=logging.DEBUG, filename="spider.log", format=format)

logger = logging.getLogger(__name__)

# ---------------------------------------------------------------------------- #
#                           Initialize browser driver                          #
# ---------------------------------------------------------------------------- #
driver = webdriver.Firefox()
driver.get("https://justjoin.it/")

# ---------------------------------------------------------------------------- #
#              Find class name for elements containing a job role              #
# ---------------------------------------------------------------------------- #
element = driver.find_element(By.XPATH, "//*[contains(text(), 'Developer')]")
job_title_element_class = element.get_attribute("class")

logger.debug(f"Determined that elements containing a job title have class {job_title_element_class}")

# Find scrollable element
scrollable_container = None
find_div_elements_wider_than_200px_xpath = "//div[(@style and contains(@style, 'width') and number(substring-before(substring-after(@style, 'width:'), 'px')) > 200) or (number(substring-before(substring-after(@width, 'px'), 'px')) > 200)]"
elements = driver.find_elements(By.XPATH, find_div_elements_wider_than_200px_xpath)
for elem in elements:
    class_name = elem.get_attribute("class")
    if "css" in class_name:
        scrollable_container = class_name
        logger.debug(f"Found scrollable container: {class_name}")

# ---------------------------------------------------------------------------- #
#                           Iterate through listings                           #
# ---------------------------------------------------------------------------- #
import time

SCROLL_PAUSE_TIME = 0.5

# Get scroll height
#while True:
for i in range(3):
    # Scrape location
    location_xpath = "//*[local-name()='svg']/*[local-name()='path' and @d='M12 12c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm6-1.8C18 6.57 15.35 4 12 4s-6 2.57-6 6.2c0 2.34 1.95 5.44 6 9.14 4.05-3.7 6-6.8 6-9.14zM12 2c4.2 0 8 3.22 8 8.2 0 3.32-2.67 7.25-8 11.8-5.33-4.55-8-8.48-8-11.8C4 5.22 7.8 2 12 2z']/../.."
    locations = driver.find_elements(By.XPATH, location_xpath)

    # Scrape company names
    company_xpath = "//*[local-name()='svg']/*[local-name()='path' and @d='M12 7V3H2v18h20V7H12zM6 19H4v-2h2v2zm0-4H4v-2h2v2zm0-4H4V9h2v2zm0-4H4V5h2v2zm4 12H8v-2h2v2zm0-4H8v-2h2v2zm0-4H8V9h2v2zm0-4H8V5h2v2zm10 12h-8v-2h2v-2h-2v-2h2v-2h-2V9h8v10zm-2-8h-2v2h2v-2zm0 4h-2v2h2v-2z']/../.."
    company_names = driver.find_elements(By.XPATH, company_xpath)

    # Scrape job titles
    job_titles = driver.find_elements(By.CLASS_NAME, job_title_element_class)

    # Scrape salaries
    salaries = driver.find_elements(By.XPATH, "//div[contains(text(), 'PLN') or contains(text(), 'Undisclosed Salary')]")
    salaries = list(filter(lambda x: x.text and x.text.strip(), salaries))

    # Scrape offer URLs
    offer_urls = driver.find_elements(By.XPATH, "//a[contains(@href, 'offers')]")

    # Print data in a tabular form
    columns = ["Location", "Company", "Job Title", "Salary", "Offer URL"]
    format_row = "{:>30}{:>50}{:>70}{:>30}{:>120}"
    print("-" * 300)
    print(format_row.format(*columns))
    print("-" * 300)
    for (location, company, job_title, salary, offer_url) in zip(locations, company_names, job_titles, salaries, offer_urls):
        location_description = location.find_elements(By.XPATH, "./span")
        location_text = " ".join(map(lambda x: x.text, location_description))

        print(format_row.format(location_text, company.text, job_title.text, salary.text, offer_url.get_attribute("href")))

    # Scroll down to the bottom
    driver.execute_script(f"document.querySelector('div.{scrollable_container}').scrollBy(0, document.body.scrollHeight);")

    # Wait to load page
    time.sleep(SCROLL_PAUSE_TIME)
