**#Overview**
##Go-Web-Crawler is a web crawler and simple search engine built in Go. It crawls specified URLs, processes the content, and stores relevant information in a SQL database. Users can search for specific words, ##and the program returns the pages (links) with the highest relevance based on their search terms.

###Features
Web Crawling: Efficiently crawls web pages, extracts words, and stores them in a SQL database.
Search Engine: Allows users to search for specific terms and retrieves the most relevant pages.
Sitemap Handling: Supports fetching and parsing sitemaps to discover new URLs for crawling.
Robots.txt Compliance: Respects website crawling policies defined in robots.txt files.
TF-IDF Calculation: Implements Term Frequency-Inverse Document Frequency (TF-IDF) to rank pages based on search terms.
Image Search: Provides functionality to search for images related to the given search terms.
