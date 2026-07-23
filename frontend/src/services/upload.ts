import axios, { type AxiosProgressEvent } from 'axios'

type ProgressEventHandler = (progressEvent: AxiosProgressEvent) => void

interface UploadService {
  uploadProductImage: (
    productId: string,
    file: File,
    backendURL: string,
    onProgress: ProgressEventHandler,
  ) => ReturnType<typeof axios.patch>
}

const upload_svc: UploadService = {
  uploadProductImage(productId, file, backendURL, onProgress) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', file.type)

    return axios.patch(`${backendURL}/api/products/${productId}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: onProgress,
    })
  },
}

export default upload_svc