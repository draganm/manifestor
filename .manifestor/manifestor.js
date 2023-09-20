render('deployment.yaml',{"name": "mfstr"}, "deployment.yaml")
render(
    'ingress.yaml', 
    {
        "hostnames": ["www.netice9.com","netice9.com"], 
        'serviceName': "foo", 
        'name': "bar",
    },
    "ingress.yaml",
)