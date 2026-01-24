#pragma once

#include <memory>

namespace oiio_go {

// Handle template for shared_ptr types (ImageCache, ColorProcessor)
// Allows Go to manage C++ shared_ptr lifecycle through opaque pointer
template<typename T>
struct Handle {
    std::shared_ptr<T> ptr;
    explicit Handle(std::shared_ptr<T> p) : ptr(std::move(p)) {}
    T* get() { return ptr.get(); }
    const T* get() const { return ptr.get(); }
};

// UniqueHandle template for unique_ptr types (ImageInput, ImageOutput)
// Allows Go to manage C++ unique_ptr lifecycle through opaque pointer
template<typename T>
struct UniqueHandle {
    std::unique_ptr<T> ptr;
    explicit UniqueHandle(std::unique_ptr<T> p) : ptr(std::move(p)) {}
    T* get() { return ptr.get(); }
    const T* get() const { return ptr.get(); }
};

} // namespace oiio_go
