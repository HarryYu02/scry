from sources.stsfandom import STSFandomScraper


def main():
    s = STSFandomScraper()
    _ = s.fetch()


if __name__ == "__main__":
    main()
