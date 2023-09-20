render('deployment.yaml',{"name": "mfstr"})
render('ingress.yaml', 
    {
        "hostnames": ["www.netice9.com","netice9.com"], 
        'serviceName': "foo", 
        'name': "bar",
    }
)