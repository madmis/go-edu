<html>
<head>
    <title>Login Form</title>
</head>
<body>
<form action="/login" method="post">
    <h2>Login Form</h2>
    <p>
        Username:<input type="text" name="username">
        Password:<input type="password" name="password">
    </p>
    <p>
        <input type="radio" name="gender" value="1">Male
        <input type="radio" name="gender" value="2">Female
    </p>
    <p>
        <input type="checkbox" name="interest" value="football">Football
        <input type="checkbox" name="interest" value="basketball">Basketball
        <input type="checkbox" name="interest" value="tennis">Tennis
    </p>
    <input type="hidden" name="token" value="{{.}}">
    <p><input type="submit" value="Login"></p>
</form>
</body>
</html>