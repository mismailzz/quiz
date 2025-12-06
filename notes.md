

---

# 🐛 Quiz Input Timeout Bug — Root Cause & Fix Summary

## ❗ The Problem

When a quiz question timed out, the program behaved incorrectly on future questions:

* The user sometimes had to press **Enter twice**
* Inputs were being **echoed multiple times**
* The quiz could become **stuck** unless additional input was entered
* Leftover goroutines kept reading from `stdin`

### **Why it happened**

For **each question**, the program spawned a new goroutine like this:

```go
go func() {
    text, _ := reader.ReadString('\n')
    answerCh <- strings.TrimSpace(text)
}()
```

If the user **did not answer before the timeout**, that goroutine:

* Stayed blocked forever in `ReadString`
* Never exited
* Stayed alive past the timeout
* Continued reading input during later questions

After several questions, multiple goroutines were all reading from `stdin` simultaneously — causing duplicated input, out-of-order input, and stuck behavior.

This situation is known as a **goroutine leak**.

---

## ✅ The Conceptual Fix

### **1. Never spawn a per-question input goroutine.**

Instead of creating a new goroutine for each question (which may never finish), the solution is to run **one single goroutine** for the entire quiz that continuously reads user input.

### **2. Send every user-entered line into a shared channel.**

```go
inputCh := make(chan string)

go func() {
    reader := bufio.NewReader(os.Stdin)
    for {
        line, _ := reader.ReadString('\n')
        inputCh <- strings.TrimSpace(line)
    }
}()
```

This goroutine:

* Never dies
* Never multiplies
* Never blocks the program
* Ensures only **one listener** touches `stdin`

### **3. Each question reads from the same input channel**

A question now simply waits for:

* The next user input from `inputCh`, **or**
* The timeout

Only one message is consumed per question, so user input always stays in sync.

---

## 🎯 End Result

* No goroutine leaks
* No duplicated input
* No stuck behavior
* Clean timeouts
* Predictable, stable quiz behavior

The quiz now behaves exactly as expected, even when the user answers late or after a timeout.

---

