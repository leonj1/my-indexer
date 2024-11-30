from setuptools import setup, find_packages

setup(
    name="my-indexer-client",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
        "pydantic>=2.0.0",
    ],
    python_requires=">=3.7",
    author="Your Name",
    description="Python client for my-indexer",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
)
