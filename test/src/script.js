async function test(){
    let response = await fetch("/auth",{
        method: "POST",
        body: "message1"
    })
    let result = await response.text()
    alert(result)
}