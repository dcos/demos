#!/usr/bin/env python
"""
Import data from StackOverflow
"""

import json
import stackexchange

def main():
    """
    Query SO and output to file
    """

    user_api_key = 'YOURSTACKOVERFLOWAPIKEYHERE'
    tags = ['mesos', 'dcos', 'mesosphere']

    # Query StackOverflow
    so_connect = stackexchange.Site(stackexchange.StackOverflow,
                                    app_key=user_api_key,
                                    impose_throttling=True)
    so_res = []
    for tag in tags:
        for question in so_connect.questions(tagged=tag, body=True):
            so_res.append(question)
        for question in so_connect.search(intitle=tag):
            full_question = so_connect.question(question.id)
            so_res.append(full_question)

    output_file = open("stack.json", "w")

    if so_res:

        # Since we used both API's we may have dupes, remove them here
        deduped = {r.id: r for r in so_res}.values()

        # Get the JSON output field
        jsonq = [d.json for d in deduped]

        for jquestion in jsonq:
            output_file.write("{}\n".format(json.dumps(jquestion)))

    output_file.close()


if __name__ == '__main__':
    main()
