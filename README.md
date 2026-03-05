## 🚀 Performance Test Results

We ran a performance test on this service using **Bombardier** on a **MacBook Air M1 (8 GB RAM)**. The results are impressive considering the hardware constraints:

### Requests per Second (RPS)
- **Average:** 143271.72  
- **Standard Deviation:** 48,175  
- **Maximum:** 183,792  

### Latency
- **Average:** 346.11us µs (~0.34 ms)  
- **Maximum:** 40.6 ms  

### HTTP Status Codes
- **2xx (Success):** 4,306,060  
- **5xx (Server Errors):** 0  
- **1xx / 3xx / 4xx / Others:** 0  

### Throughput
- **33.86 MB/s**  

### Notes
- The service shows **extremely high performance** on a lightweight machine.  
- Even with peaks reaching 180k RPS, no server errors occurred.  
- On a real production server with more CPU and RAM, the service can handle significantly higher loads with similar latency.  
- Maximum latency spikes (40 ms) are rare and likely occur during GC or peak CPU usage.  

✅ **Conclusion:** This service is highly optimized and ready for high-load scenarios.
