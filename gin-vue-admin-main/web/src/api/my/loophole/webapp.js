import request from "@/utils/request";
export const ProductDataProps = {
    name: "",
    remarks: undefined,
    provider: undefined,
    id: undefined,
};

/**
 * 获取组件列表
 * @param params
 */
export const getProductList = (params = {
    page: 1,
    pagesize: 9999,
    search_query: undefined,
}) => {
    return request({
        url: "/my/product/",
        method: "get",
        params
    });
};

/**
 * 创建组件
 * @param data
 */
export const createProduct = (ProductDataProps) => {
    return request({
        url: "/my/product/",
        method: "post",
        data:ProductDataProps
    });
};

/**
 * 删除组件
 * @param id
 */
export const deleteProduct = (id) => {
    return request({
        url: `/my/product/${id}/`,
        method: "delete"
    });
};
/**
 * 编辑组件
 * @param data
 * @param id
 */
export const updateProduct = (ProductDataProps, id) => {
    return request({
        url: `/my/product/${id}/`,
        method: "put",
        data: ProductDataProps
    });
};
