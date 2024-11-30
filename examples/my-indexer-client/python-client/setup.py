from setuptools import setup, find_packages

setup(
    name="my-indexer-client",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[],  # No external dependencies needed when using as library
    python_requires=">=3.7",
    author="Your Name",
    description="Python client library for my-indexer - A full-text search and document indexing library",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Topic :: Text Processing :: Indexing",
    ],
)
