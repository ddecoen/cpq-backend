# CPQ Backend API

🚀 **Configure, Price, Quote** backend system with enterprise licensing and AI add-ons.

## Features

✅ **Enterprise Licenses** - Starter (10), Growth (50), Scale (200), Unlimited users  
✅ **AI Add-ons** - Assistant ($15), Analytics ($25), Security ($35)  
✅ **Volume Discounts** - 20% at $50K, 30% at $100K annual value  
✅ **Multi-year Discounts** - 15% for 2+ years, 25% for 3+ years  
✅ **Customer Tier Discounts** - 10% startup, 5% growth  

## Quick Start

```bash
go mod tidy
go run .

cat >> README.md << 'EOF'

## Accessing the Application

After running the Go program, the server will start on **port 8080**. You can access the application through your web browser or API client:

### Web Interface
- **Main API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### API Endpoints
- **Products**: `GET http://localhost:8080/api/v1/demo/products`
- **Pricing**: `GET http://localhost:8080/api/v1/demo/pricing`
- **Create Quote**: `POST http://localhost:8080/api/v1/demo/quote`
- **List Quotes**: `GET http://localhost:8080/api/v1/demo/quotes`

### Testing with curl
```bash
# Check if the server is running
curl http://localhost:8080/health

# Get available products
curl http://localhost:8080/api/v1/demo/products

# Get pricing information
curl http://localhost:8080/api/v1/demo/pricing

