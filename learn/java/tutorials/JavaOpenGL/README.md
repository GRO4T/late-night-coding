# Lightweight Java OpenGL game library

## Exporting obj models from blender
File->Export->Wavefront(.obj) <br>
*Make sure to check following: Include Normals, Include UVs, Triangulate Faces*

## Calculating normals
### Finite difference method
```java
private Vector3f calculateNormal(int x, int z, BufferedImage image){
    float heightL = getHeight(x - 1, z, image);
    float heightR = getHeight(x + 1, z, image);
    float heightD = getHeight(x , z - 1, image);
    float heightU = getHeight(x, z + 1, image);
    Vector3f normal = new Vector3f(heightL - heightR, 2f, heightD - heightU);
    normal.normalise();
    return normal;
}
```

## Calculating terrain height
### How to determine height of the point inside of a triangle
#### Barycentric method
https://codeplea.com/triangular-interpolation
```java
public static float baryCentric(Vector3f p1, Vector3f p2, Vector3f p3, Vector2f pos) {
    float det = (p2.z - p3.z) * (p1.x - p3.x) + (p3.x - p2.x) * (p1.z - p3.z);
    float l1 = ((p2.z - p3.z) * (pos.x - p3.x) + (p3.x - p2.x) * (pos.y - p3.z)) / det;
    float l2 = ((p3.z - p1.z) * (pos.x - p3.x) + (p1.x - p3.x) * (pos.y - p3.z)) / det;
    float l3 = 1.0f - l1 - l2;
    return l1 * p1.y + l2 * p2.y + l3 * p3.y;
}
```