一般来说，生产系统的应用程序，应该有动态调整日志级别的功能。继续查看源码，你会发现，这个程序也可以调整日志级别。如果你给它发送 SIGUSR1 信号，就可以把日志调整为 INFO 级；发送 SIGUSR2 信号，则会调整为 WARNING 级：

```python
def set_logging_info(signal_num, frame): 
  '''Set loging level to INFO when receives SIGUSR1''' 
  logger.setLevel(logging.INFO) 

def set_logging_warning(signal_num, frame): 
  '''Set loging level to WARNING when receives SIGUSR2''' 
  logger.setLevel(logging.WARNING) 

signal.signal(signal.SIGUSR1, set_logging_info) 
signal.signal(signal.SIGUSR2, set_logging_warning) 
```

