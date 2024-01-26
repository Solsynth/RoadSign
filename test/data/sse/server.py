import time

from flask import Flask, render_template, Response

app = Flask(__name__, template_folder=".")


# Generator function to simulate real-time updates
def event_stream():
    count = 0
    while True:
        time.sleep(1)
        count += 1
        yield f"data: {count}\n\n"


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/sse')
def sse():
    return Response(event_stream(), content_type='text/event-stream')


if __name__ == '__main__':
    app.run(debug=True, threaded=True)
