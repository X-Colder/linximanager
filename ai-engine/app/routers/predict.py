"""预测 API 路由（占位实现，供部署验证）"""

from fastapi import APIRouter

router = APIRouter()


@router.post("/predict/sales", tags=["predict"])
async def predict_sales(merchant_id: int, product_id: int, days: int = 7):
    """
    销售量预测接口
    - merchant_id: 商家ID
    - product_id: 商品ID
    - days: 预测天数（默认7天）
    """
    # TODO: 接入真实预测模型
    return {
        "merchant_id": merchant_id,
        "product_id": product_id,
        "predictions": [
            {"date": f"day_{i + 1}", "qty": 10.0} for i in range(days)
        ],
    }
