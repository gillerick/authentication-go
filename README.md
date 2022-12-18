### Go Authentication Patterns

- mux's `sessions` module is used to create new sessions and validate existing ones on the fly
- In Go, sessions can be stored in the program memory by creating a `CookieStore`. The line below explicitly tells the
  program to create one by picking the secret ket from the environment variable called `SESSION_SECRET`:
  ```go
  var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
  ```
- The `NewCookieStore` function returns a store which can be used to manage cookies. If the session does not exist, an
  empty store will be returned.
  ```go
  session, _ := store.Get(r, "session.id")
  ```