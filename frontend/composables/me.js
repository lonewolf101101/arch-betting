export const useCustomer = () => {
  return useState("customer", () => null)
}

export const useFetchMe = async () => {
  const customer = useCustomer()
  const { data: c, error } = await useFetch("/api/me")
  if (error.value) {
    customer.value = null
  }
  customer.value = c.value
}
