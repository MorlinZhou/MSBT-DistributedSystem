package server;
import api.PublicService;

public class ServiceImpl implements PublicService{
    @Override
    public String sayHello(String name) {
        return "Hello, " + name + "!";
    }
}
