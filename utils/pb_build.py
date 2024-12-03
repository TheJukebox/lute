import os
import subprocess as ps
from pathlib import Path

import click


def build_proto(path: Path) -> Path | Exception:
    cmd = f"protoc -I {path.parent} --go_out=. --go-grpc_out=. {path}"
    try:
        ps.run(cmd.split(" "), check=True, stderr=ps.PIPE)
        return Path(f"gen/{path.stem}/{path.stem}.pb.go")
    except ps.CalledProcessError as e:
        click.echo(click.style(f"\nFailed to build {path}!", fg="red"))
        click.echo(click.style(f"\nprotoc error: {e.stderr.strip()}\n", fg="red"))
        return None


def find_proto_dir(d: Path =".") -> Path:
    for d, dirs, files in os.walk(d):
        if any(file.endswith(".proto") for file in files):
            return d
    return None


def build_all(path: Path) -> Path:
    proto_dir = find_proto_dir(path)
    if not proto_dir:
        raise click.exceptions.ClickException("Could not find any .proto files!")

    proto_files = []
    for d, dirs, files in os.walk(proto_dir):
        found_files = [f"{d}/{file}" for file in files if file.endswith('.proto')]
        proto_files.extend(found_files)
            
    click.echo(f"Found {len(proto_files)} .proto files: ")
    for proto in proto_files:
        click.echo(f"\t{proto}")

    for proto in proto_files:
        click.echo(f"\nBuilding: {proto}...")
        built = build_proto(Path(proto))
        if built and built.exists():
            click.echo(click.style(f"Successfully built '{built}' !", fg="green"))


@click.command()
@click.option("--all", is_flag=True, help="Compiles all protobuffers found in the current directory tree.")
@click.argument("path", required=True)
def build(path: Path, all: bool) -> None:
    if all:
        build_all(path)

if __name__ == "__main__":
    build()