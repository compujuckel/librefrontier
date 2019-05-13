<html lang="en">
<meta charset="utf-8" />
<head>
    <title>Favorites - LibreFrontier</title>
    <link rel="stylesheet" href="https://unpkg.com/chota">
</head>
<body>
<main class="container">
    <div class="row">
        <div class="col">
            <table>
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                </tr>
                </thead>
                <tbody>
                {{ range .favorites }}
                <tr>
                    <td>{{ .Id }}</td>
                    <td><a href="http://www.radio-browser.info/gui/#!/byname/{{ .Name }}">{{ .Name }}</a></td>
                    <td class="is-right"><a class="button error" href="/api/favorite/remove/{{ .Id }}">Remove</a></td>
                </tr>
                {{ end }}
                </tbody>
            </table>
        </div>
    </div>
</main>
</body>
</html>